package adverity

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/devoteamgcloud/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceAdverityLookup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"argument": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"id_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"text": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"search_terms": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"filtered_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"match_exact_term": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expect_string": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"disable_lookup": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		ReadContext: dataSourceLookupRead,
	}
}

func dataSourceLookupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	list := []string{}
	if !d.Get("disable_lookup").(bool) {
		parameters := d.Get("parameters").([]interface{})
		url := d.Get("url").(string)
		providerConfig := m.(*config)
		client := *providerConfig.Client
		params := []adverityclient.Query{}
		for _, parameter := range parameters {
			p := parameter.(map[string]interface{})
			param := adverityclient.Query{
				Key:   p["argument"].(string),
				Value: p["value"].(string),
			}
			params = append(params, param)
		}
		var res *adverityclient.LookupString
		var err error
		if d.Get("expect_string").(bool) {
			res, err = client.DoLookupString(url, params)
		} else {
			int_res, int_err := client.DoLookup(url, params)
			err = int_err
			id_mappings := []adverityclient.IDMappingString{}
			for _, id_mapping_int := range int_res.Results {
				id := strconv.Itoa(id_mapping_int.ID)
				id_mapping := adverityclient.IDMappingString{
					ID:   id,
					Name: id_mapping_int.Name,
				}
				id_mappings = append(id_mappings, id_mapping)
			}
			res = &adverityclient.LookupString{
				Error:   int_res.Error,
				Results: id_mappings,
			}
		}
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
			return diags
		}
		if res.Error != "nil" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error while doing lookup: %s", res.Error),
			})
		}
		idMappings := flattenLookup(&res.Results)
		if err := d.Set("id_mappings", idMappings); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
			return diags
		}
		search_terms, exists := d.GetOk("search_terms")
		filtered_list := []string{}
		match_exact := d.Get("match_exact_term").(bool)
		if exists {
			for _, term := range search_terms.([]interface{}) {
				if term == nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed doing lookup: empty string not permitted",
						Detail:   "In and Adverity looku an empty string (\"\") is not permitted. If you don't want to specify any search terms, leave the list empty.",
					})
					return diags
				}
				found_match := false
				for _, mapping := range idMappings {
					mapping_cast := mapping.(map[string]interface{})
					id := mapping_cast["id"].(string)
					name := mapping_cast["text"].(string)
					if name == term.(string) || (strings.Contains(strings.ToLower(name), strings.ToLower(term.(string))) && !match_exact) {
						filtered_list = append(filtered_list, id)
						found_match = true
					}
				}
				if !found_match {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Error while doing lookup: could not find a match for term \"%s\"", term.(string)),
						Detail:   fmt.Sprintf("No matches where found in the Adverity API for your search term \"%s\". If this is the first time running with this search term, double check that the term you have given is correct. If it used to work with this same search term, chances are that the search term no longer exists, or has been renamed in your platform. In this case, remove the term or replace it with its new equivalent.", term.(string)),
					})
					return diags
				}
			}
			if diags.HasError() {
				return diags
			}
		}
		allFilters := make(map[string]bool)
		for _, item := range filtered_list {
			if _, value := allFilters[item]; !value {
				allFilters[item] = true
				list = append(list, item)
			}
		}
	}
	d.Set("filtered_list", list)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func flattenLookup(idMappings *[]adverityclient.IDMappingString) []interface{} {
	if idMappings != nil {
		mappings := make([]interface{}, len(*idMappings), len(*idMappings))

		for i, idMapping := range *idMappings {
			mapping := make(map[string]interface{})
			mapping["id"] = idMapping.ID
			mapping["text"] = idMapping.Name
			mappings[i] = mapping
		}
		return mappings
	}
	return make([]interface{}, 0)
}
