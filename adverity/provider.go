package adverity

import (
	"context"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	NAME                  = "name"
	STACK                 = "stack"
	DESTINATION_TYPE      = "destination_type"
	PROJECT_ID            = "project_id"
	DATASET_ID            = "dataset_id"
	AUTH                  = "auth"
	CONNECTION_TYPE_ID    = "connection_type_id"
	CONNECTION_PARAMETERS = "connection_parameters"
	DATALAKE_ID           = "datalake_id"
	PARENT_ID             = "parent_id"
	SLUG                  = "slug"
	INSTANCE_URL          = "instance_url"
	TOKEN                 = "token"
	WORKSPACE_ID          = "workspace_id"
	SCHEMA_MAPPING        = "schema_mapping"
	CONNECTION_ID         = "connection_id"
	URL                   = "url"
	DESTINATION_ID        = "destination_id"
	DATASTREAM_ID         = "datastream_id"
	TABLE_NAME            = "table_name"
	IS_AUTHORIZED         = "is_authorized"
	HEADERS_FORMATTING    = "headers_formatting"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			INSTANCE_URL: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Url YOUR_STACK.datatap.adverity.com",
			},
			TOKEN: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"adverity_workspace":           workspace(),
			"adverity_connection":          connection(),
			"adverity_storage":             storage(),
			"adverity_destination":         destination(),
			"adverity_datastream":          datastream(),
			"adverity_destination_mapping": destinationMapping(),
			"adverity_datatype_mapping":    datatypeMapping(),
			"adverity_fetch":               fetch(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"adverity_workspace": datasourceWorkspace(),
			"adverity_auth_url":  datasourceAuthUrl(),
			"adverity_lookup":    datasourceAdverityLookup(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type config struct {
	Client *adverityclient.Client
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	client, err := adverityclient.CreateClientFromLogin(
		d.Get(INSTANCE_URL).(string),
		d.Get(TOKEN).(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	config := config{
		Client: client,
	}
	return &config, diags
}
