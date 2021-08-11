package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"example.com/adverityclient"
	"strconv"
)

func destination() *schema.Resource {
	return &schema.Resource{
		Create: destinationCreate,
		Read:   destinationRead,
		Update: destinationUpdate,
		Delete: destinationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"stack": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"destination_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dataset_id": &schema.Schema{
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



func destinationCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	stack := d.Get("stack").(int)
	auth := d.Get("auth").(int)
	destination_type := d.Get("destination_type").(int)
	project_id := d.Get("project_id").(string)
	dataset_id := d.Get("dataset_id").(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DestinationConfig{
		Name:     name,
		Stack: stack,
        ProjectID: project_id,
        DatasetID: dataset_id,
        Auth: auth,
	}

    res, err := client.CreateDestination(conf, destination_type)

	if err != nil {
		return err
	}

    d.SetId(strconv.Itoa(res.ID))
	return destinationRead(d, m)
}





func destinationRead(d *schema.ResourceData, m interface{}) error {

	destination_type := d.Get("destination_type").(int)

    providerConfig := m.(*config)

    client := *providerConfig.Client

	res, err := client.ReadDestination(d.Id(), destination_type)
	if err != nil {
		return err
	}

    d.SetId(d.Id())
	d.Set("project_id", res.Project)
	d.Set("stack", res.Stack)
	d.Set("auth", res.Auth)
	d.Set("dataset_id", res.Dataset)
	d.Set("name", res.Name)


	return nil
}

func destinationUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	stack := d.Get("stack").(int)
	auth := d.Get("auth").(int)
	destination_type := d.Get("destination_type").(int)
	project_id := d.Get("project_id").(string)
	dataset_id := d.Get("dataset_id").(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DestinationConfig{
		Name:     name,
		Stack: stack,
        ProjectID: project_id,
        DatasetID: dataset_id,
        Auth: auth,
	}

	_, err := client.UpdateDestination(conf, destination_type, d.Id())

	if err != nil {
		return err
	}
	return destinationRead(d, m)
}

func destinationDelete(d *schema.ResourceData, m interface{}) error {
    destination_type := d.Get("destination_type").(int)
	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteDestination(d.Id(), destination_type)

	if err != nil {
		return err
	}

	return nil
}
