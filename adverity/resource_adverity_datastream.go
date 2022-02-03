package adverity

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datastream() *schema.Resource {
	return &schema.Resource{
		CreateContext: datastreamCreate,
		ReadContext:   datastreamRead,
		UpdateContext: datastreamUpdate,
		DeleteContext: datastreamDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					description := val.(string)
					if len(description) > 1000 {
						errs = append(errs, fmt.Errorf("%q must be under 1000 characters, current length: %d", key, len(description)))
					}
					return
				},
			},
			"retention_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					r_type := val.(int)
					if r_type < 1 || r_type > 4 {
						errs = append(errs, fmt.Errorf("%q must be an integer between 1 and 4, got %d", key, r_type))
					}
					return
				},
				Description: "Retention Type options: 1: Retain All, 2: Retain N fetches, 3: Retain N days, 4: Retain N extracts",
			},
			"retention_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					r_number := val.(int)
					if r_number < 0 || r_number > 32767 {
						errs = append(errs, fmt.Errorf("%q must be an integer between 0 and 32767, got %d", key, r_number))
					}
					return
				},
				Description: "The amount (N) of fetches/extracts/days to retain (raw extracts are not counted). Must be an integer greater than zero.",
			},
			"overwrite_key_columns": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"overwrite_datastream": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"overwrite_filename": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_insights_mediaplan": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"manage_extract_names": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"extract_name_keys": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					name_key := val.(string)
					if len(name_key) > 128 {
						errs = append(errs, fmt.Errorf("%q must be under 128 characters, current length: %d", key, len(name_key)))
					}
					return
				},
			},
			"stack": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"auth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"datatype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"do_fetch_on_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"days_to_fetch": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
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

func datastreamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		Schedules:         &schs,
	}

	// TODO: these next few parameters use a deprecated function (GetOkExists). There is no replacement for it yet.
	// Once there is a replacement use that instead. Alternatively, find a way around the issue.
	// Issue with using GetOk: GetOk returns true for the second value if the key has been set to a non-zero value.
	// Since false for booleans and a 0 for ints is sometimes also a value we want, we can't use this method.
	if description, exists := d.GetOkExists("description"); exists {
		desc := description.(string)
		conf.Description = &desc
	}
	if retention_type, exists := d.GetOkExists("retention_type"); exists {
		ret_type := retention_type.(int)
		conf.RetentionType = &ret_type
	}
	if retention_number, exists := d.GetOkExists("retention_number"); exists {
		ret_num := retention_number.(int)
		conf.RetentionNumber = &ret_num
	}
	if overwrite_key_columns, exists := d.GetOkExists("overwrite_key_columns"); exists {
		over_key_clm := overwrite_key_columns.(bool)
		conf.OverwriteKeyColumns = &over_key_clm
	}
	if overwrite_datastream, exists := d.GetOkExists("overwrite_datastream"); exists {
		over_dtstrm := overwrite_datastream.(bool)
		conf.OverwriteDatastream = &over_dtstrm
	} else {
	}
	if overwrite_filename, exists := d.GetOkExists("overwrite_filename"); exists {
		over_filnm := overwrite_filename.(bool)
		conf.OverwriteFileName = &over_filnm
	}
	if is_insights_mediaplan, exists := d.GetOkExists("is_insights_mediaplan"); exists {
		is_ins_medplan := is_insights_mediaplan.(bool)
		conf.IsInsightsMediaplan = &is_ins_medplan
	}
	if manage_extract_names, exists := d.GetOkExists("manage_extract_names"); exists {
		mng_extract_keys := manage_extract_names.(bool)
		conf.ManageExtractNames = &mng_extract_keys
	}
	if extract_name_keys, exists := d.GetOkExists("extract_name_keys"); exists {
		extract_nm_keys := extract_name_keys.(string)
		conf.ExtractNameKeys = &extract_nm_keys
	}

	enabledConf := adverityclient.DataStreamEnablingConfig{
		Enabled: enabled,
	}

	res, err := client.CreateDatastream(conf, datastream_type_id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(res.ID))

	_, enablingErr := client.EnableDatastream(enabledConf, d.Id())

	if enablingErr != nil {
		return diag.FromErr(err)
	}

	return datastreamRead(ctx, d, m)
}

func datastreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	datastream_type_id := d.Get("datastream_type_id").(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadDatastream(d.Id(), datastream_type_id)
	if err != nil {
		return diag.FromErr(err)
	}
	schedules := flattenSchedulesData(&res.Schedules)
	if err := d.Set("schedules", schedules); err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", res.Name)
	d.Set("description", res.Description)
	d.Set("stack", res.StackID)
	d.Set("enabled", res.Enabled)
	d.Set("auth", res.Auth)
	d.Set("datatype", res.Datatype)
	d.Set("retention_type", res.RetentionType)
	d.Set("retention_number", res.RetentionNumber)
	d.Set("overwrite_key_columns", res.OverwriteKeyColumns)
	d.Set("overwrite_datastream", res.OverwriteDatastream)
	d.Set("overwrite_filename", res.OverwriteFileName)
	d.Set("is_insights_mediaplan", res.IsInsightsMediaplan)
	d.Set("manage_extract_names", res.ManageExtractNames)
	d.Set("extract_name_keys", res.ExtractNameKeys)

	return diags
}

func datastreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfig := m.(*config)
	client := *providerConfig.Client

	common_conf := adverityclient.DatastreamCommonUpdateConfig{}
	// TODO: these next few parameters use a deprecated function (GetOkExists). There is no replacement for it yet.
	// Once there is a replacement use that instead. Alternatively, find a way around the issue.
	// Issue with using GetOk: GetOk returns true for the second value if the key has been set to a non-zero value.
	// Since false for booleans and a 0 for ints is sometimes also a value we want, we can't use this method.
	if name, exists := d.GetOkExists("name"); exists {
		nm := name.(string)
		common_conf.Name = &nm
	}
	if description, exists := d.GetOkExists("description"); exists {
		desc := description.(string)
		common_conf.Description = &desc
	}
	if retention_type, exists := d.GetOkExists("retention_type"); exists {
		ret_type := retention_type.(int)
		common_conf.RetentionType = &ret_type
	}
	if retention_number, exists := d.GetOkExists("retention_number"); exists {
		ret_num := retention_number.(int)
		common_conf.RetentionNumber = &ret_num
	}
	if overwrite_key_columns, exists := d.GetOkExists("overwrite_key_columns"); exists {
		over_key_clm := overwrite_key_columns.(bool)
		common_conf.OverwriteKeyColumns = &over_key_clm
	}
	if overwrite_datastream, exists := d.GetOkExists("overwrite_datastream"); exists {
		over_dtstrm := overwrite_datastream.(bool)
		common_conf.OverwriteDatastream = &over_dtstrm
	} else {
	}
	if overwrite_filename, exists := d.GetOkExists("overwrite_filename"); exists {
		over_filnm := overwrite_filename.(bool)
		common_conf.OverwriteFileName = &over_filnm
	}
	if is_insights_mediaplan, exists := d.GetOkExists("is_insights_mediaplan"); exists {
		is_ins_medplan := is_insights_mediaplan.(bool)
		common_conf.IsInsightsMediaplan = &is_ins_medplan
	}
	if manage_extract_names, exists := d.GetOkExists("manage_extract_names"); exists {
		mng_extract_keys := manage_extract_names.(bool)
		common_conf.ManageExtractNames = &mng_extract_keys
	}
	if extract_name_keys, exists := d.GetOkExists("extract_name_keys"); exists {
		extract_nm_keys := extract_name_keys.(string)
		common_conf.ExtractNameKeys = &extract_nm_keys
	}
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

	common_conf.Schedules = schs
	_, err := client.UpdateDatastreamCommon(common_conf, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	datatype := d.Get("datatype").(string)
	datatypeConf := adverityclient.DatastreamDatatypeConfig{
		Datatype: datatype,
	}
	_, err = client.DataStreamChangeDatatype(datatypeConf, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	enabled := d.Get("enabled").(bool)
	enabledConf := adverityclient.DataStreamEnablingConfig{
		Enabled: enabled,
	}
	_, enablingErr := client.EnableDatastream(enabledConf, d.Id())
	if enablingErr != nil {
		return diag.FromErr(err)
	}

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
	specific_conf := adverityclient.DatastreamSpecificConfig{
		Parameters:        parameters,
		ParametersListInt: parameters_list_int,
		ParametersListStr: parameters_list_string,
	}
	_, err = client.UpdateDatastreamSpecific(specific_conf, d.Id(), datastream_type_id)
	if err != nil {
		return diag.FromErr(err)
	}

	if enabled {
		if d.Get("do_fetch_on_update").(bool) {
			_, err := client.ScheduleFetch(d.Get("days_to_fetch").(int), d.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return datastreamRead(ctx, d, m)
}

func datastreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	datastream_type_id := d.Get("datastream_type_id").(int)
	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteDatastream(d.Id(), datastream_type_id)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
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
