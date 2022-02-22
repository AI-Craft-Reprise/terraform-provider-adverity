package adverity

import (
	"context"
	"strconv"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func storage() *schema.Resource {
	return &schema.Resource{
		CreateContext: storageCreate,
		ReadContext:   storageRead,
		UpdateContext: storageUpdate,
		DeleteContext: storageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: storageImportHelper,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func storageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	url := d.Get("url").(string)
	stack := d.Get("stack").(int)
	auth := d.Get("auth").(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.StorageConfig{
		Name:  name,
		Stack: stack,
		Auth:  auth,
		URL:   url,
	}

	res, err := client.CreateStorage(conf)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(res.ID))

	return storageRead(ctx, d, m)
}

func storageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err, code := client.ReadStorage(d.Id())
	if err != nil {
		if code == 404 {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set("name", res.Name)
	d.Set("stack", res.Stack)
	d.Set("auth", res.Auth)
	d.Set("url", res.URL)

	return diags
}

func storageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	url := d.Get("url").(string)
	stack := d.Get("stack").(int)
	auth := d.Get("auth").(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.StorageConfig{
		Name:  name,
		Stack: stack,
		Auth:  auth,
		URL:   url,
	}

	_, err := client.UpdateStorage(conf, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}
	return storageRead(ctx, d, m)
}

func storageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteStorage(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// Useless function, included in case format changes in the future
func storageImportHelper(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	d.SetId(d.Id())

	return []*schema.ResourceData{d}, nil
}
