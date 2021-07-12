package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"datalake_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Read: workspaceDataSource,
	}
}

func workspaceDataSource(d *schema.ResourceData, m interface{}) error {

    workspace_id := d.Get("workspace_id").(string)
    slug := d.Get("slug").(string)

    providerConfig := m.(*config)

    client := *providerConfig.Client


	res, err := client.ReadWorkspace(slug)
	if err != nil {
		return err
	}
    d.SetId(workspace_id)
	d.Set("datalake_id", res.Datalake)
	d.Set("parent_id", res.ParentID)
	d.Set("name", res.Name)
    d.Set("slug", res.Slug)

	return nil
}
