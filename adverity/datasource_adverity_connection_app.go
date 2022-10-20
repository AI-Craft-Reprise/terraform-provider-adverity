package adverity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceAdverityConnectionApp() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"connection_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The connection type ID for the connection for which you're looking up the app ID.",
			},
			"selector": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional selector for when there is the possibility of finding multiple apps. Make sure this corresponds to the name of the app in the API.",
			},
			"app": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the app or the connection type.",
			},
		},
		ReadContext: datasourceConnectionApp,
		Description: "An app ID lookup for connections that are authorised through an app (as opposed to with service account keys or other methods).",
	}
}

func datasourceConnectionApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connectionTypeID := d.Get("connection_type_id").(int)
	selector := d.Get("selector").(string)
	providerConfig := m.(*config)
	client := *providerConfig.Client
	result, err := client.LookupConnectionApp(connectionTypeID, selector)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("app", result)
	d.SetId(fmt.Sprintf("%d/%d", connectionTypeID, result))
	return diags
}
