package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"example.com/adverityclient"
	"strconv"
	"log"
)

func workspace() *schema.Resource {
	return &schema.Resource{
		Create: workspaceCreate,
		Read:   workspaceRead,
		Update: workspaceUpdate,
		Delete: workspaceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"datalake_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}



func workspaceCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	datalake_id := d.Get("datalake_id").(string)
	parent_id := d.Get("parent_id").(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.CreateWorkspaceConfig{
		Name:     name,
		DatalakeID:    datalake_id,
		ParentID:    parent_id,
	}

    res, err := client.CreateWorkspace(conf)

	if err != nil {
		return err
	}

    d.Set("slug", res.Slug)

	return workspaceRead(d, m)
}





func workspaceRead(d *schema.ResourceData, m interface{}) error {

    slug := d.Get("slug").(string)
    datalake_id := d.Get("datalake_id").(string)
    parent_id := d.Get("parent_id").(int)
    name := d.Get("name").(string)

    providerConfig := m.(*config)

    client := *providerConfig.Client


	res, err := client.ReadWorkspace(slug)
	if err != nil {
		return err
	}
    d.SetId(strconv.Itoa(res.ID))
	d.Set("datalake_id", datalake_id)
	d.Set("parent_id", parent_id)
	d.Set("name", name)


	return nil
}

func workspaceUpdate(d *schema.ResourceData, m interface{}) error {
	parent_id := d.Get("parent_id").(int)
    name := d.Get("name").(string)
    slug := d.Get("slug").(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.UpdateWorkspaceConfig{
	    Name:     name,
		StackSlug:     slug,
		ParentID:    parent_id,

	}

	_, err := client.UpdateWorkspace(conf)

	if err != nil {
		return err
	}
	return workspaceRead(d, m)
}

func workspaceDelete(d *schema.ResourceData, m interface{}) error {

	providerConfig := m.(*config)


	client := *providerConfig.Client
    log.Println(d.Get("slug").(string))
	conf := adverityclient.DeleteWorkspaceConfig{
		StackSlug: d.Get("slug").(string),
	}

	_, err := client.DeleteWorkspace(conf)

	if err != nil {
		return err
	}

	return nil
}
