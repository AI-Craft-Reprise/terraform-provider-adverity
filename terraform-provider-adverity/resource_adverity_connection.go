package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"example.com/adverityclient"
	"strconv"
)


func connection() *schema.Resource {
	return &schema.Resource{
		Create: connectionCreate,
		Read:   connectionRead,
		Update: connectionUpdate,
		Delete: connectionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"stack": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"connection_type_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"connection_parameters": {
                Type:     schema.TypeMap,
                Optional: true,
                Elem: &schema.Schema{
                        Type: schema.TypeString,
                },
            },
		},
	}
}



func connectionCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	stack := d.Get("stack").(int)
	connection_type_id := d.Get("connection_type_id").(int)

    connection_parameters, exists := d.GetOk("connection_parameters")

    parameters:=[]*adverityclient.ConnectionParameters{}

	if exists {
		for n, v := range connection_parameters.(map[string]interface{}) {
			parameter:=new(adverityclient.ConnectionParameters)
			parameter.Value=v.(string)
			parameter.Name=n
			parameters=append(parameters,parameter)
		}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.ConnectionConfig{
		Name:     name,
		Stack:    stack,
		ConnectionParameters: parameters,
	}

    res, err := client.CreateConnection(conf, connection_type_id)

	if err != nil {
		return err
	}

    d.SetId(strconv.Itoa(res.ID))
	//TODO here we should authorized the application or not

	return connectionRead(d, m)
}





func connectionRead(d *schema.ResourceData, m interface{}) error {

	connection_type_id := d.Get("connection_type_id").(int)

    providerConfig := m.(*config)

    client := *providerConfig.Client


	res, err := client.ReadConnection(d.Id(),connection_type_id)
	if err != nil {
		return err
	}
	d.Set("name", res.Name)
	d.Set("stack", res.Stack)


	return nil
}

func connectionUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	stack := d.Get("stack").(int)
	connection_type_id := d.Get("connection_type_id").(int)

	connection_parameters, exists := d.GetOk("connection_parameters")

    parameters:=[]*adverityclient.ConnectionParameters{}
    if exists {
		for n, v := range connection_parameters.(map[string]interface{}) {
			parameter:=new(adverityclient.ConnectionParameters)
			parameter.Value=v.(string)
			parameter.Name=n
			parameters=append(parameters,parameter)
		}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.ConnectionConfig{
		Name:     name,
		Stack:    stack,
		ConnectionParameters: parameters,
	}

	_, err := client.UpdateConnection(conf, d.Id(),connection_type_id)

	if err != nil {
		return err
	}
	return connectionRead(d, m)
}

func connectionDelete(d *schema.ResourceData, m interface{}) error {
    connection_type_id := d.Get("connection_type_id").(int)
	providerConfig := m.(*config)


	client := *providerConfig.Client

	_, err := client.DeleteConnection(d.Id(), connection_type_id)

	if err != nil {
		return err
	}

	return nil
}
