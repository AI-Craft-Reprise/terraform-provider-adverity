package adverity

import (
	"context"
	"strconv"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// 	"log"
)

func workspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: workspaceCreate,
		ReadContext:   workspaceRead,
		UpdateContext: workspaceUpdate,
		DeleteContext: workspaceDelete,

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

func workspaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	d.Set(SLUG, res.Slug)

	return workspaceRead(ctx, d, m)
}

func workspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	slug := d.Get(SLUG).(string)
	datalakeId := d.Get(DATALAKE_ID).(string)
	parentId := d.Get(PARENT_ID).(int)
	name := d.Get(NAME).(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadWorkspace(slug)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(res.ID))
	d.Set(DATALAKE_ID, datalakeId)
	d.Set(PARENT_ID, parentId)
	d.Set(NAME, name)

	return diags
}

func workspaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	return workspaceRead(ctx, d, m)
}

func workspaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfig := m.(*config)

	client := *providerConfig.Client
	conf := adverityclient.DeleteWorkspaceConfig{
		StackSlug: d.Get(SLUG).(string),
	}

	_, err := client.DeleteWorkspace(conf)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
