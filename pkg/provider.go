//
// DISCLAIMER
//
// Copyright 2020-2022 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//

package pkg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				DefaultFunc: schema.EnvDefaultFunc("OASIS_ENDPOINT", "api.cloud.arangodb.com"),
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
			"oasis_deployment":                   resourceDeployment(),
			"oasis_ipallowlist":                  resourceIPAllowlist(),
			"oasis_certificate":                  resourceCertificate(),
			"oasis_backup":                       resourceBackup(),
			"oasis_backup_policy":                resourceBackupPolicy(),
			"oasis_project":                      resourceProject(),
			"oasis_example_dataset_installation": resourceExampleDatasetInstallation(),
			"oasis_organization":                 resourceOrganization(),
			"oasis_iam_group":                    resourceIAMGroup(),
			"oasis_iam_role":                     resourceIAMRole(),
			"oasis_organization_invite":          resourceOrganizationInvite(),
			"oasis_auditlog":                     resourceAuditLog(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"oasis_project":                       dataSourceOasisProject(),
			"oasis_organization":                  dataSourceOasisOrganization(),
			"oasis_terms_and_conditions":          dataSourceTermsAndConditions(),
			"oasis_example_dataset_installations": dataSourceOasisExampleDatasetInstallation(),
			"oasis_example_datasets":              dataSourceOasisExampleDataset(),
			"oasis_backup":                        dataSourceOasisBackup(),
			"oasis_cloud_provider":                dataSourceOasisCloudProvider(),
			"oasis_region":                        dataSourceOasisRegion(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Initialize Client with connection settings
	client := Client{
		ApiKeyID:      d.Get("api_key_id").(string),
		ApiKeySecret:  d.Get("api_key_secret").(string),
		ApiEndpoint:   d.Get("oasis_endpoint").(string),
		ApiPortSuffix: d.Get("api_port_suffix").(string),
	}
	if v, ok := d.GetOk("project"); ok {
		client.ProjectID = v.(string)
	}
	if v, ok := d.GetOk("organization"); ok {
		client.OrganizationID = v.(string)
	}
	return &client, nil
}
