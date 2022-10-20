package adverity

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceAuthUrl() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			CONNECTION_TYPE_ID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The connection type ID for the connection the auth url belongs to.",
			},
			CONNECTION_ID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the connection the auth url belongs to.",
			},
			URL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url to authorise the connection.",
			},
		},
		ReadContext: authUrlDataSource,
		Description: "This datasource will generate an authentication url for a connection. This url, when followed, will authenticate the connection. The url will change everytime this datasource is run.",
	}
}

func authUrlDataSource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connection_type_id := d.Get(CONNECTION_TYPE_ID).(string)
	connection_id := d.Get(CONNECTION_ID).(string)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err := client.ReadAuthUrl(connection_type_id, connection_id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(CONNECTION_TYPE_ID, connection_type_id)
	d.Set(CONNECTION_ID, connection_id)
	d.Set(URL, res.URL)

	return diags
}
