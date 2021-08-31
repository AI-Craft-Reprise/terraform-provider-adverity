package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		Read: workspaceDataSource,
	}
}

func workspaceDataSource(d *schema.ResourceData, m interface{}) error {

	workspace_id := d.Get(WORKSPACE_ID).(string)
	slug := d.Get(SLUG).(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadWorkspace(slug)
	if err != nil {
		return err
	}
	d.SetId(workspace_id)
	d.Set(DATALAKE_ID, res.Datalake)
	d.Set(PARENT_ID, res.ParentID)
	d.Set(NAME, res.Name)
	d.Set(SLUG, res.Slug)

	return nil
}
