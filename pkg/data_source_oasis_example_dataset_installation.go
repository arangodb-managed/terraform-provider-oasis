//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
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
// Author Gergely Brautigam
//

package pkg

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	example "github.com/arangodb-managed/apis/example/v1"
)

var (
	resourceName                            = "dataset-installation"
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
		Read: dataSourceOasisExampleDatasetInstallationRead,

		Schema: map[string]*schema.Schema{
			installationDeploymentIdFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			installationItemsFieldName: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						installationExampleDatasetIdFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						installationCreatedAtFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						installationStatusFieldName: {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									installationStatusDatabaseNameFieldName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									installationStatusStateFieldName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									installationStatusIsAvailableFieldName: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									installationStatusIsFailedFieldName: {
										Type:     schema.TypeBool,
										Computed: true,
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
func dataSourceOasisExampleDatasetInstallationRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	examplec := example.NewExampleDatasetServiceClient(client.conn)
	deplID := data.Get(installationDeploymentIdFieldName).(string)
	response, err := examplec.ListExampleDatasetInstallations(client.ctxWithToken, &example.ListExampleDatasetInstallationsRequest{
		DeploymentId: deplID,
	})
	if err != nil {
		client.log.Error().Str("deployment-id", deplID).Err(err).Msg("Failed to get list of example installations for deployment.")
		return err
	}

	for k, v := range flattenExampleDatasetInstallations(deplID, response.GetItems()) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(uniqueResourceID(resourceName))
	return nil
}

// flattenExampleDatasetInstallations takes a list of installations and converts them into a terraform consumable format.
func flattenExampleDatasetInstallations(id string, items []*example.ExampleDatasetInstallation) map[string]interface{} {
	return map[string]interface{}{
		installationDeploymentIdFieldName: id,
		installationItemsFieldName:        flattenExampleDatasetInstallation(items),
	}
}

// flattenExampleDatasetInstallation takes a single installation and converts it into a terraform consumable format.
func flattenExampleDatasetInstallation(items []*example.ExampleDatasetInstallation) []interface{} {
	ret := make([]interface{}, 0)
	for _, v := range items {
		ret = append(ret, map[string]interface{}{
			installationExampleDatasetIdFieldName: v.GetExampledatasetId(),
			installationCreatedAtFieldName:        v.GetCreatedAt().String(),
			installationStatusFieldName:           flattenStatus(v.GetStatus()),
		})
	}
	return ret
}

// flattenStatus takes the status portion of the installation and converts it into a terraform consumable format.
func flattenStatus(status *example.ExampleDatasetInstallation_Status) []interface{} {
	return []interface{}{
		map[string]interface{}{
			installationStatusDatabaseNameFieldName: status.GetDatabaseName(),
			installationStatusStateFieldName:        status.GetState(),
			installationStatusIsFailedFieldName:     status.GetIsFailed(),
			installationStatusIsAvailableFieldName:  status.GetIsAvailable(),
		},
	}
}
