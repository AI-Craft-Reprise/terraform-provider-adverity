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
				Type:     schema.TypeInt,
				Required: true,
			},
			"app": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		ReadContext: datasourceConnectionApp,
	}
}

func datasourceConnectionApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connectionTypeID := d.Get("connection_type_id").(int)
	providerConfig := m.(*config)
	client := *providerConfig.Client
	result, err := client.LookupConnectionApp(connectionTypeID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("app", result)
	d.SetId(fmt.Sprintf("%d/%d", connectionTypeID, result))
	return diags
}
