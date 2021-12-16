package adverity

import (
	"strconv"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func destination() *schema.Resource {
	return &schema.Resource{
		Create: destinationCreate,
		Read:   destinationRead,
		Update: destinationUpdate,
		Delete: destinationDelete,

		Schema: map[string]*schema.Schema{
			NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			STACK: {
				Type:     schema.TypeInt,
				Required: true,
			},
			DESTINATION_TYPE: {
				Type:     schema.TypeInt,
				Required: true,
			},
			PROJECT_ID: {
				Type:     schema.TypeString,
				Required: true,
			},
			DATASET_ID: {
				Type:     schema.TypeString,
				Required: true,
			},
			AUTH: {
				Type:     schema.TypeInt,
				Required: true,
			},
			SCHEMA_MAPPING: {
				Type:     schema.TypeBool,
				Required: true,
			},
			HEADERS_FORMATTING: {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func destinationCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get(NAME).(string)
	stack := d.Get(STACK).(int)
	auth := d.Get(AUTH).(int)
	destinationType := d.Get(DESTINATION_TYPE).(int)
	projectId := d.Get(PROJECT_ID).(string)
	datasetId := d.Get(DATASET_ID).(string)
	schemaMapping := d.Get(SCHEMA_MAPPING).(bool)
	headersFormatting := d.Get(HEADERS_FORMATTING).(int)

	if headersFormatting <= 0 || headersFormatting > 3 {
		return errorString{"Could not create Destination. Invalid value " + strconv.Itoa(headersFormatting) + " for headers_formatting. Only 1, 2, or 3 is allowed."}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DestinationConfig{
		Name:              name,
		Stack:             stack,
		ProjectID:         projectId,
		DatasetID:         datasetId,
		Auth:              auth,
		SchemaMapping:     schemaMapping,
		HeadersFormatting: headersFormatting,
	}

	res, err := client.CreateDestination(conf, destinationType)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(res.ID))
	return destinationRead(d, m)
}

func destinationRead(d *schema.ResourceData, m interface{}) error {

	destinationType := d.Get(DESTINATION_TYPE).(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadDestination(d.Id(), destinationType)
	if err != nil {
		return err
	}

	d.SetId(d.Id())
	d.Set(PROJECT_ID, res.Project)
	d.Set(STACK, res.Stack)
	d.Set(AUTH, res.Auth)
	d.Set(DATASET_ID, res.Dataset)
	d.Set(NAME, res.Name)
	d.Set(SCHEMA_MAPPING, res.SchemaMapping)
	d.Set(HEADERS_FORMATTING, res.HeadersFormatting)

	return nil
}

func destinationUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get(NAME).(string)
	stack := d.Get(STACK).(int)
	auth := d.Get(AUTH).(int)
	destinationType := d.Get(DESTINATION_TYPE).(int)
	projectId := d.Get(PROJECT_ID).(string)
	datasetId := d.Get(DATASET_ID).(string)
	schemaMapping := d.Get(SCHEMA_MAPPING).(bool)
	headersFormatting := d.Get(HEADERS_FORMATTING).(int)

	if headersFormatting <= 0 || headersFormatting > 3 {
		return errorString{"Could not create Destination. Invalid value " + strconv.Itoa(headersFormatting) + " for headers_formatting. Only 1, 2, or 3 is allowed."}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.DestinationConfig{
		Name:              name,
		Stack:             stack,
		ProjectID:         projectId,
		DatasetID:         datasetId,
		Auth:              auth,
		SchemaMapping:     schemaMapping,
		HeadersFormatting: headersFormatting,
	}

	_, err := client.UpdateDestination(conf, destinationType, d.Id())

	if err != nil {
		return err
	}
	return destinationRead(d, m)
}

func destinationDelete(d *schema.ResourceData, m interface{}) error {
	destinationType := d.Get(DESTINATION_TYPE).(int)
	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteDestination(d.Id(), destinationType)

	if err != nil {
		return err
	}

	return nil
}
