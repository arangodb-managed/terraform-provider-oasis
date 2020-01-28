//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Joerg Schad
//

package pkg

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider defines an ArangoDB Oasis Terraform provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OASIS_API_KEY_ID", ""),
				Description: "OASIS API KEY ID",
			},
			"api_key_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OASIS_API_KEY_SECRET", ""),
				Description: "OASIS API KEY SECRET",
			},
			"oasis_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OASIS_ENDPOINT", "api.cloud.adbtest.xyz"),
				Description: "OASIS API ENDPOINT",
			},
			"api_port_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OASIS_PORT_SUFFIX", ":443"),
				Description: "OASIS API PORT SUFFIX",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OASIS_ORGANIZATION", ""),
				Description: "Default Oasis Organization",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OASIS_PROJECT", ""),
				Description: "Default Oasis Project",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"oasis_deployment": resourceDeployment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"oasis_organization": dataSourceOasisOrganization(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Initialize Client with connection settings
	client := Client{
		ApiKeyID:      d.Get("api_key_id").(string),
		ApiKeySecret:  d.Get("api_key_secret").(string),
		ApiEndpoint:   d.Get("oasis_endpoint").(string),
		ApiPortSuffix: d.Get("api_port_suffix").(string),
	}
	return &client, nil
}
