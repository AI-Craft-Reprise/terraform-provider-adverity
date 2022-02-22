package adverity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NAME: {
				Type:     schema.TypeString,
				Computed: true,
			},
			DATALAKE_ID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			PARENT_ID: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			SLUG: {
				Type:     schema.TypeString,
				Required: true,
			},
			WORKSPACE_ID: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ReadContext: workspaceDataSource,
	}
}

func workspaceDataSource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	workspace_id := d.Get(WORKSPACE_ID).(string)
	slug := d.Get(SLUG).(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err, code := client.ReadWorkspace(slug)
	if err != nil {
		if code == 404 {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	d.SetId(workspace_id)
	d.Set(DATALAKE_ID, res.Datalake)
	d.Set(PARENT_ID, res.ParentID)
	d.Set(NAME, res.Name)
	d.Set(SLUG, res.Slug)

	return diags
}
