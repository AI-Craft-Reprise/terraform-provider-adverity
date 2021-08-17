package main

import (
	"example.com/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
)

func connection() *schema.Resource {
	return &schema.Resource{
		Create: connectionCreate,
		Read:   connectionRead,
		Update: connectionUpdate,
		Delete: connectionDelete,

		Schema: map[string]*schema.Schema{
			NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			STACK: {
				Type:     schema.TypeInt,
				Required: true,
			},
			CONNECTION_TYPE_ID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			CONNECTION_PARAMETERS: {
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
	name := d.Get(NAME).(string)
	stack := d.Get(STACK).(int)
	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)

	connectionParameters, exists := d.GetOk(CONNECTION_PARAMETERS)

	parameters := []*adverityclient.Parameters{}

	if exists {
		for n, v := range connectionParameters.(map[string]interface{}) {
			parameter := new(adverityclient.Parameters)
			parameter.Value = v.(string)
			parameter.Name = n
			parameters = append(parameters, parameter)
		}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.ConnectionConfig{
		Name:       name,
		Stack:      stack,
		Parameters: parameters,
	}

	res, err := client.CreateConnection(conf, connectionTypeId)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(res.ID))
	//TODO here we should authorized the application or not

	return connectionRead(d, m)
}

func connectionRead(d *schema.ResourceData, m interface{}) error {

	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadConnection(d.Id(), connectionTypeId)
	if err != nil {
		return err
	}
	d.Set(NAME, res.Name)
	d.Set(STACK, res.Stack)

	return nil
}

func connectionUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get(NAME).(string)
	stack := d.Get(STACK).(int)
	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)

	connectionParameters, exists := d.GetOk(CONNECTION_PARAMETERS)

	parameters := []*adverityclient.Parameters{}
	if exists {
		for n, v := range connectionParameters.(map[string]interface{}) {
			parameter := new(adverityclient.Parameters)
			parameter.Value = v.(string)
			parameter.Name = n
			parameters = append(parameters, parameter)
		}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.ConnectionConfig{
		Name:       name,
		Stack:      stack,
		Parameters: parameters,
	}

	_, err := client.UpdateConnection(conf, d.Id(), connectionTypeId)

	if err != nil {
		return err
	}
	return connectionRead(d, m)
}

func connectionDelete(d *schema.ResourceData, m interface{}) error {
	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)
	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteConnection(d.Id(), connectionTypeId)

	if err != nil {
		return err
	}

	return nil
}
