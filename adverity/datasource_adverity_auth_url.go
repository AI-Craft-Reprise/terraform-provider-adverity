package adverity

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceAuthUrl() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			CONNECTION_TYPE_ID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			CONNECTION_ID: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			AUTH_URL: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: authUrlDataSource,
	}
}

func authUrlDataSource(d *schema.ResourceData, m interface{}) error {
	connection_type_id := d.Get(CONNECTION_TYPE_ID).(int)
	connection_id := d.Get(CONNECTION_ID).(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadAuthUrl(connection_type_id, connection_id)
	if err != nil {
		return err
	}
	d.Set(CONNECTION_TYPE_ID, connection_type_id)
	d.Set(CONNECTION_ID, connection_id)
	d.Set(AUTH_URL, res.URL)

	return nil
}
