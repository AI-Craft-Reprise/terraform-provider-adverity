package adverity

import (
	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	// 	"log"
)

func workspace() *schema.Resource {
	return &schema.Resource{
		Create: workspaceCreate,
		Read:   workspaceRead,
		Update: workspaceUpdate,
		Delete: workspaceDelete,

		Schema: map[string]*schema.Schema{
			NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			DATALAKE_ID: {
				Type:     schema.TypeString,
				Required: true,
			},
			PARENT_ID: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			SLUG: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func workspaceCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get(NAME).(string)
	datalakeId := d.Get(DATALAKE_ID).(string)
	parentId := d.Get(PARENT_ID).(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.CreateWorkspaceConfig{
		Name:       name,
		DatalakeID: datalakeId,
		ParentID:   parentId,
	}

	res, err := client.CreateWorkspace(conf)

	if err != nil {
		return err
	}

	d.Set(SLUG, res.Slug)

	return workspaceRead(d, m)
}

func workspaceRead(d *schema.ResourceData, m interface{}) error {

	slug := d.Get(SLUG).(string)
	datalakeId := d.Get(DATALAKE_ID).(string)
	parentId := d.Get(PARENT_ID).(int)
	name := d.Get(NAME).(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadWorkspace(slug)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(res.ID))
	d.Set(DATALAKE_ID, datalakeId)
	d.Set(PARENT_ID, parentId)
	d.Set(NAME, name)

	return nil
}

func workspaceUpdate(d *schema.ResourceData, m interface{}) error {
	parentId := d.Get(PARENT_ID).(int)
	name := d.Get(NAME).(string)
	slug := d.Get(SLUG).(string)
	datalake_id := d.Get("datalake_id").(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.UpdateWorkspaceConfig{
		Name:       name,
		StackSlug:  slug,
		ParentID:   parentId,
		DatalakeID: datalake_id,
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
	//log.Println(d.Get(SLUG).(string))
	conf := adverityclient.DeleteWorkspaceConfig{
		StackSlug: d.Get(SLUG).(string),
	}

	_, err := client.DeleteWorkspace(conf)

	if err != nil {
		return err
	}

	return nil
}
