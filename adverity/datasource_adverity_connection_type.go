package adverity

import (
	"context"
	"fmt"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceAdverityConnectionType() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_search_term": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug_search_term": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_type_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		ReadContext: datasourceConnectionType,
	}
}

func datasourceConnectionType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	searchTerm := d.Get("api_search_term").(string)
	providerConfig := m.(*config)
	client := *providerConfig.Client
	results, err := client.LookupConnectionTypes(searchTerm)
	if err != nil {
		return diag.FromErr(err)
	}
	finalResults := []adverityclient.ConnectionType{}
	if len(results) <= 0 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("No results were found for search term %s", searchTerm),
		})
	} else if len(results) > 1 {
		slugSearchTerm := d.Get("slug_search_term").(string)
		for _, result := range results {
			if result.Slug == slugSearchTerm {
				finalResults = append(finalResults, result)
			}
		}
		if len(finalResults) <= 0 {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("No results were found for slug search term %s", slugSearchTerm),
			})
		} else if len(finalResults) > 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Multiple results (%d) were found for search term %s and slug search term %s, the first result will be selected (%s)", len(finalResults), searchTerm, slugSearchTerm, finalResults[0].Name),
			})
		}
	} else {
		finalResults = results
	}
	d.Set("connection_type_id", finalResults[0].ID)
	id, err := uuid.GenerateUUID()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(id)
	return diags
}
