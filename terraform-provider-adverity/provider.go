package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"example.com/adverityclient"
)


func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"instance_url": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Url YOUR_STACK.datatap.adverity.com",
			},
			"token": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"adverity_workspace":             workspace(),
			"adverity_connection":             connection(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"adverity_workspace": datasourceWorkspace(),
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
	client, err := adverityclient.CreateClientFromLogin(d.Get("instance_url").(string),
		d.Get("token").(string))

	if err != nil {
		return nil, err
	}

	config := config{
		Client: client,
	}

	return &config, nil

}
