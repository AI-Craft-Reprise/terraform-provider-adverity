package adverity

import (
	"strconv"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func destinationMapping() *schema.Resource {
	return &schema.Resource{
		Create: destinationMappingCreate,
		Read:   destinationMappingRead,
		Update: destinationMappingUpdate,
		Delete: destinationMappingDelete,

		Schema: map[string]*schema.Schema{
			DESTINATION_TYPE: {
				Type:     schema.TypeInt,
				Required: true,
			},
			DESTINATION_ID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			DATASTREAM_ID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			TABLE_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func destinationMappingCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}

	d.SetId(strconv.Itoa(res.ID))
	return destinationMappingRead(d, m)
}

func destinationMappingRead(d *schema.ResourceData, m interface{}) error {
	destination_type := d.Get(DESTINATION_TYPE).(int)
	destination_id := d.Get(DESTINATION_ID).(int)
	id, strErr := strconv.Atoi(d.Id())
	if strErr != nil {
		return strErr
	}
	providerConfig := m.(*config)
	client := *providerConfig.Client
	res, err := client.ReadDestinationMapping(id, destination_type, destination_id)
	if err != nil {
		return err
	}

	d.SetId(d.Id())
	d.Set(DESTINATION_TYPE, destination_type)
	d.Set(DESTINATION_ID, res.Target)
	d.Set(DATASTREAM_ID, res.Datastream)
	d.Set(TABLE_NAME, res.TableName)

	return nil
}

func destinationMappingUpdate(d *schema.ResourceData, m interface{}) error {
	destination_type := d.Get(DESTINATION_TYPE).(int)
	destination_id := d.Get(DESTINATION_ID).(int)
	datastream_id := d.Get(DATASTREAM_ID).(int)
	table_name := d.Get(TABLE_NAME).(string)
	id, strErr := strconv.Atoi(d.Id())
	if strErr != nil {
		return strErr
	}

	providerConfig := m.(*config)
	client := *providerConfig.Client

	conf := adverityclient.DestinationMappingConfig{
		Datastream: datastream_id,
		TableName:  table_name,
	}

	_, err := client.UpdateDestinationMapping(conf, destination_type, destination_id, id)

	if err != nil {
		return err
	}
	return destinationMappingRead(d, m)
}

func destinationMappingDelete(d *schema.ResourceData, m interface{}) error {
	destination_type := d.Get(DESTINATION_TYPE).(int)
	destination_id := d.Get(DESTINATION_ID).(int)
	id, strErr := strconv.Atoi(d.Id())
	if strErr != nil {
		return strErr
	}

	providerConfig := m.(*config)
	client := *providerConfig.Client

	_, err := client.DeleteDestinationMapping(id, destination_type, destination_id)

	return err
}
