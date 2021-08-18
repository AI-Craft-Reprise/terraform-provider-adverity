package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	// 	"log"
	// 	"reflect"
)

func datastream() *schema.Resource {
	return &schema.Resource{
		Create: datastreamCreate,
		Read:   datastreamRead,
		Update: datastreamUpdate,
		Delete: datastreamDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"stack": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"datastream_type_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"datastream_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"datastream_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func datastreamCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	stack := d.Get("stack").(int)
	datastream_type_id := d.Get("datastream_type_id").(int)

	datastream_parameters, exists := d.GetOk("datastream_parameters")

	parameters := []*adverityclient.Parameters{}

	if exists {
		for n, v := range datastream_parameters.(map[string]interface{}) {
			parameter := new(adverityclient.Parameters)
			parameter.Value = v.(string)
			parameter.Name = n
			parameters = append(parameters, parameter)
		}
	}

	datastream_list_map, exists := d.GetOk("datastream_list")

	parameters_list_int := []*adverityclient.ParametersListInt{}
	if exists {
		datastreamParamSet := datastream_list_map.(*schema.Set).List()
		for _, datastreamParam := range datastreamParamSet {

			for _, dp := range datastreamParam.(map[string]interface{}) {
				for _, param := range dp.([]interface{}) {
					parameter := new(adverityclient.ParametersListInt)
					values := param.(map[string]interface{})["values"]
					name := param.(map[string]interface{})["name"]
					parameter.Name = name.(string)
					for _, value := range values.([]interface{}) {
						parameter.Value = append(parameter.Value, value.(int))
					}
					parameters_list_int = append(parameters_list_int, parameter)
				}
			}

		}

	}
	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DatastreamConfig{
		Name:              name,
		Stack:             stack,
		Parameters:        parameters,
		ParametersListInt: parameters_list_int,
	}

	res, err := client.CreateDatastream(conf, datastream_type_id)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(res.ID))

	return datastreamRead(d, m)
}

func datastreamRead(d *schema.ResourceData, m interface{}) error {

	datastream_type_id := d.Get("datastream_type_id").(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadDatastream(d.Id(), datastream_type_id)
	if err != nil {
		return err
	}
	d.Set("name", res.Name)
	d.Set("stack", res.StackID)

	return nil
}

func datastreamUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	stack := d.Get("stack").(int)
	datastream_type_id := d.Get("datastream_type_id").(int)

	datastream_parameters, exists := d.GetOk("datastream_parameters")

	parameters := []*adverityclient.Parameters{}

	if exists {
		for n, v := range datastream_parameters.(map[string]interface{}) {
			parameter := new(adverityclient.Parameters)
			parameter.Value = v.(string)
			parameter.Name = n
			parameters = append(parameters, parameter)
		}
	}

	datastream_list_map, exists := d.GetOk("datastream_list")

	parameters_list_int := []*adverityclient.ParametersListInt{}
	if exists {
		datastreamParamSet := datastream_list_map.(*schema.Set).List()
		for _, datastreamParam := range datastreamParamSet {

			for _, dp := range datastreamParam.(map[string]interface{}) {
				for _, param := range dp.([]interface{}) {
					parameter := new(adverityclient.ParametersListInt)
					values := param.(map[string]interface{})["values"]
					name := param.(map[string]interface{})["name"]
					parameter.Name = name.(string)
					for _, value := range values.([]interface{}) {
						parameter.Value = append(parameter.Value, value.(int))
					}
					parameters_list_int = append(parameters_list_int, parameter)
				}
			}

		}

	}
	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DatastreamConfig{
		Name:              name,
		Stack:             stack,
		Parameters:        parameters,
		ParametersListInt: parameters_list_int,
	}

	_, err := client.UpdateDatastream(conf, d.Id(), datastream_type_id)

	if err != nil {
		return err
	}
	return datastreamRead(d, m)
}

func datastreamDelete(d *schema.ResourceData, m interface{}) error {
	datastream_type_id := d.Get("datastream_type_id").(int)
	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteDatastream(d.Id(), datastream_type_id)

	if err != nil {
		return err
	}

	return nil
}
