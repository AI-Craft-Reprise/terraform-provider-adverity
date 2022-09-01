package adverity

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/devoteamgcloud/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func columns() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the datastream.",
			},
			"schema": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				Description: "A JSON schema with all the required columns and their datatype.",
			},
			"ignored_columns": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "A list of columns which may be present in the schema JSON, but that will be ignored when setting columns in Adverity.",
			},
		},
		CreateContext: columnsCreate,
		ReadContext:   columnsRead,
		DeleteContext: columnsDelete,
		UpdateContext: columnsUpdate,
		Description:   "This resource will manage all the colmuns and their datatypes for a given Adverity datastream. Due to the way the API works, any columns created outside of this resource for the datastream will be overridden, including deleting existing columns.",
	}
}

func columnsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	datastreamID := d.Get("datastream_id").(string)
	providerConfig := m.(*config)
	client := *providerConfig.Client
	var definedSchema []adverityclient.SchemaElementNoMode
	schemaText := d.Get("schema").(string)
	if err := json.Unmarshal([]byte(schemaText), &definedSchema); err != nil {
		return diag.FromErr(err)
	}
	var columnConfigs []adverityclient.ColumnConfig
	ignoredColumns := []string{}
	if _, exists := d.GetOk("ignored_columns"); exists {
		for _, column := range d.Get("ignored_columns").([]interface{}) {
			ignoredColumns = append(ignoredColumns, column.(string))
		}
	}
	for _, column := range definedSchema {
		toIgnore := false
		for _, ignored_column := range ignoredColumns {
			if ignored_column == column.Name {
				toIgnore = true
				break
			}
		}
		if !toIgnore {
			columnConfigs = append(columnConfigs, adverityclient.ColumnConfig{
				Name: column.Name,
				Type: typeMapping[column.Type],
			})
		}
	}
	createdColumns, err := client.CreateColumns(datastreamID, columnConfigs)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, column := range createdColumns {
		if !column.ConfirmedType {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("After creating the columns, %s has confirmedType set to false. This should not happen and is either a bug in the API or the provider.", column.Name),
			})
		}
	}
	if diags.HasError() {
		return diags
	}
	d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))
	diags = append(diags, columnsRead(ctx, d, m)...)
	return diags
}

func columnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	// Check if datastream still exists. Someone might have deleted it manually, which causes the columns to also disappear
	exists, err := client.DatastreamExists(datastreamID)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		d.SetId("")
		return diags
	}
	columns, err := client.ReadColumns(datastreamID)
	if err != nil {
		return diag.FromErr(err)
	}
	var definedSchema []adverityclient.SchemaElementNoMode
	schemaText := d.Get("schema").(string)
	if err := json.Unmarshal([]byte(schemaText), &definedSchema); err != nil {
		return diag.FromErr(err)
	}
	ignoredColumns := []string{}
	if _, exists := d.GetOk("ignored_columns"); exists {
		for _, column := range d.Get("ignored_columns").([]interface{}) {
			ignoredColumns = append(ignoredColumns, column.(string))
		}
	}
	var APISchema []adverityclient.SchemaElementNoMode
	// For every column defined in the input schema
	for _, definedColumn := range definedSchema {
		found := false
		// Go over every column that has been read from the Adverity API
		for idx, column := range columns {
			// Check if the names match
			if definedColumn.Name == column.Name {
				// Add the column read from the API to the read schema
				// If the column  has been schema mapped (has a target column), add mapped = true
				APISchema = append(APISchema, adverityclient.SchemaElementNoMode{
					Type:   typeMapping[column.DataType],
					Name:   column.Name,
					Mapped: column.TargetColumn != nil,
				})
				// Remove that column from the list of columns read from the API
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
				if definedColumn.Name == ignoredColumn {
					// Add it to the read schema
					APISchema = append(APISchema, definedColumn)
					break
				}
			}
		}
	}
	// For every remaining column that has been read in the API (and thus had no match in the input schema)
	for _, column := range columns {
		// Add it to the read schema
		// If the column  has been schema mapped (has a target column), add mapped = true
		APISchema = append(APISchema, adverityclient.SchemaElementNoMode{
			Type:   typeMapping[column.DataType],
			Name:   column.Name,
			Mapped: column.TargetColumn != nil,
		})
	}
	bytes, _ := json.Marshal(APISchema)
	d.Set("schema", string(bytes[:]))
	return diags
}

func columnsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	datastreamID := d.Get("datastream_id").(string)
	providerConfig := m.(*config)
	client := *providerConfig.Client
	// Using CreateColumns with an empty slice will remove all the columns
	result, err := client.CreateColumns(datastreamID, []adverityclient.ColumnConfig{})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(result) > 0 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "After sending a POST request with an empty slice, some results got returned.",
		})
	}
	d.SetId("")
	return diags
}

func columnsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	datastreamID := d.Get("datastream_id").(string)
	providerConfig := m.(*config)
	client := *providerConfig.Client
	var definedSchema []adverityclient.SchemaElementNoMode
	schemaText := d.Get("schema").(string)
	if err := json.Unmarshal([]byte(schemaText), &definedSchema); err != nil {
		return diag.FromErr(err)
	}
	var columnConfigs []adverityclient.ColumnConfig
	ignoredColumns := []string{}
	if _, exists := d.GetOk("ignored_columns"); exists {
		for _, column := range d.Get("ignored_columns").([]interface{}) {
			ignoredColumns = append(ignoredColumns, column.(string))
		}
	}
	for _, column := range definedSchema {
		toIgnore := false
		for _, ignored_column := range ignoredColumns {
			if ignored_column == column.Name {
				toIgnore = true
				break
			}
		}
		if !toIgnore {
			columnConfigs = append(columnConfigs, adverityclient.ColumnConfig{
				Name: column.Name,
				Type: typeMapping[column.Type],
			})
		}
	}
	createdColumns, err := client.CreateColumns(datastreamID, columnConfigs)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, column := range createdColumns {
		if !column.ConfirmedType {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("After updating the columns, %s has confirmedType set to false. This should not happen and is either a bug in the API or the provider.", column.Name),
			})
		}
	}
	if diags.HasError() {
		return diags
	}
	diags = append(diags, columnsRead(ctx, d, m)...)
	return diags
}
