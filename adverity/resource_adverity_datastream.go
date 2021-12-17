package adverity

import (
	"strconv"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"auth": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"datatype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"datastream_type_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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
			"datastream_string_list": {
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
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"schedules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron_preset": {
							Type:     schema.TypeString,
							Required: true,
						},
						"time_range_preset": {
							Type:     schema.TypeInt,
							Required: true,
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
	auth := d.Get("auth").(int)
	datatype := d.Get("datatype").(string)
	datastream_type_id := d.Get("datastream_type_id").(int)
	enabled := d.Get("enabled").(bool)

	schedules := d.Get("schedules").([]interface{})
	schs := []adverityclient.Schedule{}

	for _, schedule := range schedules {
		s := schedule.(map[string]interface{})
		sch := adverityclient.Schedule{
			CronPreset:      s["cron_preset"].(string),
			TimeRangePreset: s["time_range_preset"].(int),
		}
		schs = append(schs, sch)
	}

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

	datastream_list_string_map, exists := d.GetOk("datastream_string_list")

	parameters_list_string := []*adverityclient.ParametersListStr{}
	if exists {
		datastreamParamSet := datastream_list_string_map.(*schema.Set).List()
		for _, datastreamParam := range datastreamParamSet {

			for _, dp := range datastreamParam.(map[string]interface{}) {
				for _, param := range dp.([]interface{}) {
					parameter := new(adverityclient.ParametersListStr)
					values := param.(map[string]interface{})["values"]
					name := param.(map[string]interface{})["name"]
					parameter.Name = name.(string)
					for _, value := range values.([]interface{}) {
						parameter.Value = append(parameter.Value, value.(string))
					}
					parameters_list_string = append(parameters_list_string, parameter)
				}
			}

		}

	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DatastreamConfig{
		Name:              name,
		Stack:             stack,
		Auth:              auth,
		Datatype:          datatype,
		Parameters:        parameters,
		ParametersListInt: parameters_list_int,
		ParametersListStr: parameters_list_string,
		Schedules:         schs,
	}

	enabledConf := adverityclient.DataStreamEnablingConfig{
		Enabled: enabled,
	}

	res, err := client.CreateDatastream(conf, datastream_type_id)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(res.ID))

	_, enablingErr := client.EnableDatastream(enabledConf, d.Id())

	if enablingErr != nil {
		return err
	}

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
	schedules := flattenSchedulesData(&res.Schedules)
	if err := d.Set("schedules", schedules); err != nil {
		return err
	}
	d.Set("name", res.Name)
	d.Set("stack", res.StackID)
	d.Set("enabled", res.Enabled)
	d.Set("auth", res.Auth)
	d.Set("datatype", res.Datatype)

	return nil
}

func datastreamUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	enabled := d.Get("enabled").(bool)

	schedules := d.Get("schedules").([]interface{})
	schs := []adverityclient.Schedule{}

	for _, schedule := range schedules {
		s := schedule.(map[string]interface{})
		sch := adverityclient.Schedule{
			CronPreset:      s["cron_preset"].(string),
			TimeRangePreset: s["time_range_preset"].(int),
		}
		schs = append(schs, sch)
	}

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

	datastream_list_string_map, exists := d.GetOk("datastream_string_list")

	parameters_list_string := []*adverityclient.ParametersListStr{}
	if exists {
		datastreamParamSet := datastream_list_string_map.(*schema.Set).List()
		for _, datastreamParam := range datastreamParamSet {

			for _, dp := range datastreamParam.(map[string]interface{}) {
				for _, param := range dp.([]interface{}) {
					parameter := new(adverityclient.ParametersListStr)
					values := param.(map[string]interface{})["values"]
					name := param.(map[string]interface{})["name"]
					parameter.Name = name.(string)
					for _, value := range values.([]interface{}) {
						parameter.Value = append(parameter.Value, value.(string))
					}
					parameters_list_string = append(parameters_list_string, parameter)
				}
			}

		}

	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DatastreamConfig{
		Name:              name,
		Parameters:        parameters,
		ParametersListInt: parameters_list_int,
		ParametersListStr: parameters_list_string,
		Schedules:         schs,
	}

	enabledConf := adverityclient.DataStreamEnablingConfig{
		Enabled: enabled,
	}

	datatype := d.Get("datatype").(string)

	datatypeConf := adverityclient.DatastreamDatatypeConfig{
		Datatype: datatype,
	}

	_, err := client.UpdateDatastream(conf, d.Id())

	if err != nil {
		return err
	}

	_, err = client.DataStreamChanegDatatype(datatypeConf, d.Id())

	if err != nil {
		return err
	}

	_, enablingErr := client.EnableDatastream(enabledConf, d.Id())

	if enablingErr != nil {
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

func flattenSchedulesData(schedules *[]adverityclient.Schedule) []interface{} {
	if schedules != nil {
		schs := make([]interface{}, len(*schedules), len(*schedules))

		for i, schedule := range *schedules {
			sch := make(map[string]interface{})
			sch["cron_preset"] = schedule.CronPreset
			sch["time_range_preset"] = schedule.TimeRangePreset
			schs[i] = sch
		}
		return schs
	}
	return make([]interface{}, 0)
}
