package adverity

import (
	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
)

func storage() *schema.Resource {
	return &schema.Resource{
		Create: storageCreate,
		Read:   storageRead,
		Update: storageUpdate,
		Delete: storageDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"stack": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auth": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func storageCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}

	d.SetId(strconv.Itoa(res.ID))

	return storageRead(d, m)
}

func storageRead(d *schema.ResourceData, m interface{}) error {

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadStorage(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", res.Name)
	d.Set("stack", res.Stack)
	d.Set("auth", res.Auth)
	d.Set("url", res.URL)

	return nil
}

func storageUpdate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	return connectionRead(d, m)
}

func storageDelete(d *schema.ResourceData, m interface{}) error {

	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteStorage(d.Id())

	if err != nil {
		return err
	}

	return nil
}
