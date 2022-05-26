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

	example "github.com/arangodb-managed/apis/example/v1"
)

var (
	installationResourceName                = "dataset-installation"
	installationDeploymentIdFieldName       = "deployment_id"
	installationItemsFieldName              = "items"
	installationExampleDatasetIdFieldName   = "example_dataset_id"
	installationCreatedAtFieldName          = "created_at"
	installationStatusFieldName             = "status"
	installationStatusDatabaseNameFieldName = "database_name"
	installationStatusStateFieldName        = "state"
	installationStatusIsAvailableFieldName  = "is_available"
	installationStatusIsFailedFieldName     = "is_failed"
)

// dataSourceOasisExampleDatasetInstallation defines an Example Dataset Installation datasource terraform type.
func dataSourceOasisExampleDatasetInstallation() *schema.Resource {
	return &schema.Resource{
		Description: "Example DataSet Installation Data Source",

		ReadContext: dataSourceOasisExampleDatasetInstallationRead,

		Schema: map[string]*schema.Schema{
			installationDeploymentIdFieldName: {
				Type:        schema.TypeString,
				Description: "Example Dataset Data Source Deployment ID field",
				Required:    true,
			},
			installationItemsFieldName: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						installationExampleDatasetIdFieldName: {
							Type:        schema.TypeString,
							Description: "Example Dataset Data Source ID field",
							Computed:    true,
						},
						installationCreatedAtFieldName: {
							Type:        schema.TypeString,
							Description: "Example Dataset Data Source Created At field",
							Computed:    true,
						},
						installationStatusFieldName: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									installationStatusDatabaseNameFieldName: {
										Type:        schema.TypeString,
										Description: "Example Dataset Installation Database field",
										Computed:    true,
									},
									installationStatusStateFieldName: {
										Type:        schema.TypeString,
										Description: "Example Dataset Installation State field",
										Computed:    true,
									},
									installationStatusIsAvailableFieldName: {
										Type:        schema.TypeBool,
										Description: "Example Dataset Installation IsAvailable field",
										Computed:    true,
									},
									installationStatusIsFailedFieldName: {
										Type:        schema.TypeBool,
										Description: "Example Dataset Installation Failed field",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// dataSourceOasisExampleDatasetInstallationRead reloads the resource object from the terraform store.
func dataSourceOasisExampleDatasetInstallationRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	examplec := example.NewExampleDatasetServiceClient(client.conn)
	deplID := data.Get(installationDeploymentIdFieldName).(string)
	response, err := examplec.ListExampleDatasetInstallations(client.ctxWithToken, &example.ListExampleDatasetInstallationsRequest{
		DeploymentId: deplID,
	})
	if err != nil {
		client.log.Error().Str("deployment-id", deplID).Err(err).Msg("Failed to get list of example installations for deployment.")
		return diag.FromErr(err)
	}

	for k, v := range flattenExampleDatasetInstallations(deplID, response.GetItems()) {
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	data.SetId(uniqueResourceID(installationResourceName))
	return nil
}

// flattenExampleDatasetInstallations takes result of installations and converts them into a terraform consumable format.
func flattenExampleDatasetInstallations(id string, items []*example.ExampleDatasetInstallation) map[string]interface{} {
	return map[string]interface{}{
		installationDeploymentIdFieldName: id,
		installationItemsFieldName:        flattenInstallation(items),
	}
}

// flattenExampleDatasetInstallation converts the list of installations it into a terraform consumable format.
func flattenInstallation(items []*example.ExampleDatasetInstallation) []interface{} {
	ret := make([]interface{}, 0)
	for _, v := range items {
		ret = append(ret, map[string]interface{}{
			installationExampleDatasetIdFieldName: v.GetExampledatasetId(),
			installationCreatedAtFieldName:        v.GetCreatedAt().String(),
			installationStatusFieldName:           flattenInstallationStatus(v.GetStatus()),
		})
	}
	return ret
}

// flattenStatus takes the status portion of the installation and converts it into a terraform consumable format.
func flattenInstallationStatus(status *example.ExampleDatasetInstallation_Status) []interface{} {
	return []interface{}{
		map[string]interface{}{
			installationStatusDatabaseNameFieldName: status.GetDatabaseName(),
			installationStatusStateFieldName:        status.GetState(),
			installationStatusIsFailedFieldName:     status.GetIsFailed(),
			installationStatusIsAvailableFieldName:  status.GetIsAvailable(),
		},
	}
}
