package adverity

import (
	"context"
	"strconv"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func destinationMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: destinationMappingCreate,
		ReadContext:   destinationMappingRead,
		UpdateContext: destinationMappingUpdate,
		DeleteContext: destinationMappingDelete,

		Schema: map[string]*schema.Schema{
			DESTINATION_TYPE: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			DESTINATION_ID: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			DATASTREAM_ID: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			TABLE_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastream_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"fetch_on_creation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"days_to_fetch": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
	}
}

func destinationMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	destination_type := d.Get(DESTINATION_TYPE).(int)
	destination_id := d.Get(DESTINATION_ID).(int)
	datastream_id := d.Get(DATASTREAM_ID).(int)
	table_name := d.Get(TABLE_NAME).(string)

	providerConfig := m.(*config)
	client := *providerConfig.Client

	conf := adverityclient.DestinationMappingConfig{
		Datastream: datastream_id,
		TableName:  table_name,
	}

	res, err := client.CreateDestinationMapping(conf, destination_type, destination_id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(res.ID))

	if d.Get("datastream_enabled").(bool) {
		if d.Get("fetch_on_creation").(bool) {
			_, err := client.ScheduleFetch(d.Get("days_to_fetch").(int), strconv.Itoa(datastream_id))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return destinationMappingRead(ctx, d, m)
}

func destinationMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	destination_type := d.Get(DESTINATION_TYPE).(int)
	destination_id := d.Get(DESTINATION_ID).(int)
	id, strErr := strconv.Atoi(d.Id())
	if strErr != nil {
		return diag.FromErr(strErr)
	}
	providerConfig := m.(*config)
	client := *providerConfig.Client
	res, err := client.ReadDestinationMapping(id, destination_type, destination_id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	d.Set(DESTINATION_TYPE, destination_type)
	d.Set(DESTINATION_ID, res.Target)
	d.Set(DATASTREAM_ID, res.Datastream)
	d.Set(TABLE_NAME, res.TableName)

	return diags
}

func destinationMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	destination_type := d.Get(DESTINATION_TYPE).(int)
	destination_id := d.Get(DESTINATION_ID).(int)
	datastream_id := d.Get(DATASTREAM_ID).(int)
	table_name := d.Get(TABLE_NAME).(string)
	id, strErr := strconv.Atoi(d.Id())
	if strErr != nil {
		return diag.FromErr(strErr)
	}

	providerConfig := m.(*config)
	client := *providerConfig.Client

	conf := adverityclient.DestinationMappingConfig{
		Datastream: datastream_id,
		TableName:  table_name,
	}

	_, err := client.UpdateDestinationMapping(conf, destination_type, destination_id, id)

	if err != nil {
		return diag.FromErr(err)
	}
	return destinationMappingRead(ctx, d, m)
}

func destinationMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	destination_type := d.Get(DESTINATION_TYPE).(int)
	destination_id := d.Get(DESTINATION_ID).(int)
	id, strErr := strconv.Atoi(d.Id())
	if strErr != nil {
		return diag.FromErr(strErr)
	}

	providerConfig := m.(*config)
	client := *providerConfig.Client

	_, err := client.DeleteDestinationMapping(id, destination_type, destination_id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
