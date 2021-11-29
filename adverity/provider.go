package adverity

import (
	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			INSTANCE_URL: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Url YOUR_STACK.datatap.adverity.com",
			},
			TOKEN: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"adverity_workspace":   workspace(),
			"adverity_connection":  connection(),
			"adverity_storage":     storage(),
			"adverity_destination": destination(),
			"adverity_datastream":  datastream(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"adverity_workspace": datasourceWorkspace(),
			"adverity_auth_url":  datasourceAuthUrl(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, provider, terraformVersion)
	}

	return provider
}

type config struct {
	Client *adverityclient.Client
}

func providerConfigure(d *schema.ResourceData, p *schema.Provider, terraformVersion string) (interface{}, error) {
	client, err := adverityclient.CreateClientFromLogin(
		d.Get(INSTANCE_URL).(string),
		d.Get(TOKEN).(string))

	if err != nil {
		return nil, err
	}

	config := config{
		Client: client,
	}

	return &config, nil

}
