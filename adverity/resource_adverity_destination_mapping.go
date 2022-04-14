package adverity

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
		Importer: &schema.ResourceImporter{
			StateContext: destinationMappingImportHelper,
		},

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
		},
	}
}

func destinationMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if !d.Get("datastream_enabled").(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "WARNING: Destination mapping is blocked by the datastream not being enabled, so it will not be created.",
		})
		d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))
	} else {
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
		diags = append(diags, destinationMappingRead(ctx, d, m)...)
	}
	return diags
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
	res, err, code := client.ReadDestinationMapping(id, destination_type, destination_id)
	if err != nil {
		if code == 404 {
			d.SetId("")
			return diags
		}
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

func destinationMappingImportHelper(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected destination_type:destination_id:destinationmapping_id", d.Id())
	}
	destination_type, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("could not convert destination_type (%s) to an integer", parts[0])
	}
	destination_id, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("could not convert destination_id (%s) to an integer", parts[1])
	}
	d.Set(DESTINATION_TYPE, destination_type)
	d.Set(DESTINATION_ID, destination_id)
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
