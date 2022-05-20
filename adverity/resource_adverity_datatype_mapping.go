package adverity

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func datatypeMapping() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the datastream.",
				ForceNew:    true,
			},
			"schema": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				StateFunc: func(v interface{}) string {
					jsonString, _ := structure.NormalizeJsonString(v)
					var modeElements []SchemaElementMode
					err := json.Unmarshal([]byte(jsonString), &modeElements)
					if err != nil {
						return ""
					}
					var noModeElements []SchemaElementNoMode
					for _, modeElement := range modeElements {
						noModeElements = append(noModeElements, SchemaElementNoMode{
							Type: modeElement.Type,
							Name: modeElement.Name,
						})
					}
					bytes, _ := json.Marshal(noModeElements)
					return string(bytes[:])
				},
				Description: "A JSON schema, as extracted from a BigQuery table.",
			},
			"populating_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_authorised": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"days_to_fetch": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  14,
						},
					},
				},
			},
			"mapped": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_on_missing_columns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the resource will throw an error if a column in the schema is not found in Adverity or vice versa.",
			},
			"ignored_columns": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"wait_for_columns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to true, the resource will wait until at least one column exists in the API before proceeding.",
			},
			"replace_special_characters": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to true, special characters in Adverity will be replaced by underscores, and names beginning with a number will start with an \"n\" instead.",
			},
		},
		CreateContext: datatypeMappingCreate,
		ReadContext:   datatypeMappingRead,
		UpdateContext: datatypeMappingUpdate,
		DeleteContext: datatypeMappingDelete,
		Description:   "This is an experimental resource meant to be used when exporting data without using the schema mapping feature to BigQuery. It takes as input the desired Schema that the BigQuery table has (the same format as when you would export a schema from a BigQuery table). It will read all columns associated with the specified datastream and make sure they are all the same datatype as specified in the schema. It has quite a few nested loops, so even though effort has been made to optimise it, large schemas may cause this resource to take some time. Patch requests to change columns will only be done if they are different from the desired datatype or if they are not fixed yet.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

type SchemaElementMode struct {
	Mode string
	Name string
	Type string
}

type SchemaElementNoMode struct {
	Name string
	Type string
}

func datatypeMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	typeMapping := map[string]string{
		"String":   "STRING",
		"Long":     "INTEGER",
		"Float":    "FLOAT",
		"Date":     "DATE",
		"DateTime": "DATETIME",
		"Boolean":  "BOOLEAN",
		"JSON":     "JSON",
	}
	datastreamID := d.Get("datastream_id").(string)
	providerConfig := m.(*config)
	client := *providerConfig.Client
	columns, err := client.ReadColumns(datastreamID)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("replace_special_characters").(bool) {
		columns = replaceSpecialCharacters(columns)
	}
	var existingSchema []SchemaElementNoMode
	schemaText := d.Get("schema").(string)
	if err := json.Unmarshal([]byte(schemaText), &existingSchema); err != nil {
		return diag.FromErr(err)
	}
	ignoredColumns := []string{}
	if _, exists := d.GetOk("ignored_columns"); exists {
		for _, column := range d.Get("ignored_columns").([]interface{}) {
			ignoredColumns = append(ignoredColumns, column.(string))
		}
	}
	var readSchema []SchemaElementNoMode
	// For every column defined in the input schema
	for _, existingColumn := range existingSchema {
		found := false
		// Go over every column that has been read from the Adverity API
		for idx, column := range columns {
			// Check if the names match
			if existingColumn.Name == column.Name {
				// Add the column read from the API to the read schema
				readSchema = append(readSchema, SchemaElementNoMode{
					Type: typeMapping[column.DataType],
					Name: column.Name,
				})
				// Remove that column from the mist of columns read from the API
				columns = append(columns[0:idx], columns[idx+1:]...)
				found = true
				break
			}
		}
		// If the column wasn't found in the Adverity API, it may still be one of the ignored columns. If this is the case, it has to be
		// added to the schema regardless, otherwise Terraform will (correctly) detect drift, and will try to change it. Since it is an ignored
		// column, we don't need/want Terraform to do any changes to these columns.
		if !found {
			// Go over all columns in the ignored list
			for _, ignoredColumn := range ignoredColumns {
				// If the column defined in the input schema matches a column name in the ignored list
				if existingColumn.Name == ignoredColumn {
					// Add it to the read schema
					readSchema = append(readSchema, existingColumn)
					break
				}
			}
		}
	}
	// For every remaining column that has been read in the API (and thus had no match in the input schema)
	for _, column := range columns {
		// Add it to the read schema
		readSchema = append(readSchema, SchemaElementNoMode{
			Type: typeMapping[column.DataType],
			Name: column.Name,
		})
	}
	bytes, _ := json.Marshal(readSchema)
	d.Set("schema", string(bytes[:]))
	return diags
}

func datatypeMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	typeMapping := map[string]string{
		"STRING":   "String",
		"INTEGER":  "Long",
		"FLOAT":    "Float",
		"DATE":     "Date",
		"DATETIME": "DateTime",
		"BOOLEAN":  "Boolean",
		"JSON":     "JSON",
	}
	schemaText := d.Get("schema").(string)
	var schema []SchemaElementNoMode
	if err := json.Unmarshal([]byte(schemaText), &schema); err != nil {
		return diag.FromErr(err)
	}
	providerConfig := m.(*config)
	client := *providerConfig.Client
	datastreamID := d.Get("datastream_id").(string)
	if d.Get("populating_settings.0.connection_authorised").(bool) {
		columns, err := client.ReadColumns(datastreamID)
		if d.Get("wait_for_columns").(bool) && len(columns) == 0 {
			for len(columns) == 0 {
				time.Sleep(10 * time.Second)
				columns, err = client.ReadColumns(datastreamID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("replace_special_characters").(bool) {
			columns = replaceSpecialCharacters(columns)
		}
		// notFoundInAPI and notFoundInSchema are switched naming wise, should be other way around
		notFoundInAPI := []string{}
		ignoredColumns := []string{}
		if _, exists := d.GetOk("ignored_columns"); exists {
			for _, column := range d.Get("ignored_columns").([]interface{}) {
				ignoredColumns = append(ignoredColumns, column.(string))
			}
		}
		for _, column := range columns {
			found := false
			for idx, targetColumn := range schema {
				if column.Name == targetColumn.Name {
					found = true
					if !column.ConfirmedType || column.DataType != typeMapping[targetColumn.Type] {
						client.PatchColumn(strconv.Itoa(column.ID), typeMapping[targetColumn.Type])
						log.Printf("[DEBUG] Patch request goes here to change %s with type %s to type %s", column.Name, column.DataType, typeMapping[targetColumn.Type])
					}
					// Remove found item from list to search
					schema = append(schema[0:idx], schema[idx+1:]...)
					break
				}
			}
			if !found {
				toIgnore := false
				for _, ignoredColumn := range ignoredColumns {
					if column.Name == ignoredColumn {
						toIgnore = true
					}
				}
				if !toIgnore {
					notFoundInAPI = append(notFoundInAPI, column.Name)
				}
			}
		}
		if len(schema) > 0 {
			notFoundInSchema := []string{}
			for _, column := range schema {
				toIgnore := false
				for _, ignoredColumn := range ignoredColumns {
					if column.Name == ignoredColumn {
						toIgnore = true
					}
				}
				if !toIgnore {
					notFoundInSchema = append(notFoundInSchema, column.Name)
				}
			}
			if !d.Get("error_on_missing_columns").(bool) && len(notFoundInSchema) > 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Could not find references in Adverity API for following columns specified in the schema: %s", strings.Join(notFoundInSchema, ", ")),
				})
			} else if len(notFoundInSchema) > 0 {
				return diag.Errorf("Could not find references in Adverity API for following columns specified in the schema: %s", strings.Join(notFoundInSchema, ", "))
			}
		}
		if len(notFoundInAPI) > 0 {
			if !d.Get("error_on_missing_columns").(bool) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Found references in Adverity API which are not present in the specified schema: %s", strings.Join(notFoundInAPI, ", ")),
				})
			} else {
				return diag.Errorf("Found references in Adverity API which are not present in the specified schema: %s", strings.Join(notFoundInAPI, ", "))
			}
		}
		d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))

		diags = append(diags, datatypeMappingRead(ctx, d, m)...)
		d.Set("mapped", true)
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "WARNING: Connection is not authorised yet, so it is not possible to populate the columns to map them.",
		})
		d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))
		d.Set("mapped", false)
	}

	return diags
}

func datatypeMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	typeMapping := map[string]string{
		"STRING":   "String",
		"INTEGER":  "Long",
		"FLOAT":    "Float",
		"DATE":     "Date",
		"DATETIME": "DateTime",
		"BOOLEAN":  "Boolean",
		"JSON":     "JSON",
	}
	schemaText := d.Get("schema").(string)
	var schema []SchemaElementNoMode
	if err := json.Unmarshal([]byte(schemaText), &schema); err != nil {
		return diag.FromErr(err)
	}
	providerConfig := m.(*config)
	client := *providerConfig.Client
	datastreamID := d.Get("datastream_id").(string)
	if d.Get("populating_settings.0.connection_authorised").(bool) {
		columns, err := client.ReadColumns(datastreamID)
		if d.Get("wait_for_columns").(bool) && len(columns) == 0 {
			for len(columns) == 0 {
				time.Sleep(10 * time.Second)
				columns, err = client.ReadColumns(datastreamID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("replace_special_characters").(bool) {
			columns = replaceSpecialCharacters(columns)
		}
		notFoundInAPI := []string{}
		for _, column := range columns {
			found := false
			for idx, targetColumn := range schema {
				if column.Name == targetColumn.Name {
					found = true
					if !column.ConfirmedType || column.DataType != typeMapping[targetColumn.Type] {
						client.PatchColumn(strconv.Itoa(column.ID), typeMapping[targetColumn.Type])
						log.Printf("[DEBUG] Patch request goes here to change %s with type %s to type %s", column.Name, column.DataType, typeMapping[targetColumn.Type])
					}
					// Remove found item from list to search
					schema = append(schema[0:idx], schema[idx+1:]...)
					break
				}
			}
			if !found {
				notFoundInAPI = append(notFoundInAPI, column.Name)
			}
		}
		if len(schema) > 0 {
			notFoundInSchema := []string{}
			for _, column := range schema {
				notFoundInSchema = append(notFoundInSchema, column.Name)
			}
			if !d.Get("error_on_missing_columns").(bool) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Could not find references in Adverity API for following columns specified in the schema: %s", strings.Join(notFoundInSchema, ", ")),
				})
			} else {
				return diag.Errorf("Could not find references in Adverity API for following columns specified in the schema: %s", strings.Join(notFoundInSchema, ", "))
			}
		}
		if len(notFoundInAPI) > 0 {
			if !d.Get("error_on_missing_columns").(bool) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Found references in Adverity API which are not present in the specified schema: %s", strings.Join(notFoundInAPI, ", ")),
				})
			} else {
				return diag.Errorf("Found references in Adverity API which are not present in the specified schema: %s", strings.Join(notFoundInAPI, ", "))
			}
		}
		diags = append(diags, datatypeMappingRead(ctx, d, m)...)
		d.Set("mapped", true)
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "WARNING: Connection is not authorised yet, so it is not possible to populate the columns to map them.",
		})
		d.Set("mapped", false)
	}
	return diags
}

func datatypeMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

func replaceSpecialCharacters(columns []adverityclient.Column) []adverityclient.Column {
	replacedColumns := []adverityclient.Column{}
	for _, column := range columns {
		column.Name = strings.ReplaceAll(column.Name, "(", "_")
		column.Name = strings.ReplaceAll(column.Name, ")", "_")
		column.Name = strings.ReplaceAll(column.Name, "%", "_")
		column.Name = strings.ReplaceAll(column.Name, "-", "_")
		column.Name = strings.ReplaceAll(column.Name, " ", "_")
		column.Name = strings.ReplaceAll(column.Name, ".", "_")
		column.Name = strings.ReplaceAll(column.Name, ":", "_")
		if unicode.IsDigit(rune(column.Name[0])) {
			column.Name = "n" + column.Name
		}
		replacedColumns = append(replacedColumns, column)
	}
	return replacedColumns
}
