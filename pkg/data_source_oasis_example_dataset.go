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
	exampleResourceName                        = "example-dataset"
	exampleOrganizationIDFieldName             = "organization"
	exampleExampleDatasetsFieldName            = "example_datasets"
	exampleExampleDatasetsIDFieldName          = "id"
	exampleExampleDatasetsNameFieldName        = "name"
	exampleExampleDatasetsDescriptionFieldName = "description"
	exampleExampleDatasetsGuideFieldName       = "guide"
	exampleExampleDatasetsCreatedAtFieldName   = "created_at"
)

// dataSourceOasisExampleDataset defines an Example Dataset datasource terraform type.
func dataSourceOasisExampleDataset() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOasisExampleDatasetRead,

		Schema: map[string]*schema.Schema{
			exampleOrganizationIDFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			exampleExampleDatasetsFieldName: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						exampleExampleDatasetsIDFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						exampleExampleDatasetsNameFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						exampleExampleDatasetsDescriptionFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						exampleExampleDatasetsCreatedAtFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						exampleExampleDatasetsGuideFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceOasisExampleDatasetRead reloads the resource object from the terraform store.
func dataSourceOasisExampleDatasetRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	examplec := example.NewExampleDatasetServiceClient(client.conn)
	orgID := client.OrganizationID
	if v, ok := data.GetOk(exampleOrganizationIDFieldName); ok {
		orgID = v.(string)
	}
	response, err := examplec.ListExampleDatasets(client.ctxWithToken, &example.ListExampleDatasetsRequest{
		OrganizationId: orgID,
	})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to get list of example datasets.")
		return err
	}

	for k, v := range flattenExampleDatasets(orgID, response.GetItems()) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(uniqueResourceID(exampleResourceName))
	return nil
}

// flattenExampleDatasets takes result of datasets and converts them into a terraform consumable format.
func flattenExampleDatasets(id string, items []*example.ExampleDataset) map[string]interface{} {
	return map[string]interface{}{
		exampleOrganizationIDFieldName:  id,
		exampleExampleDatasetsFieldName: flattenExampleDataset(items),
	}
}

// flattenExampleDataset converts the list of datasets it into a terraform consumable format.
func flattenExampleDataset(items []*example.ExampleDataset) []interface{} {
	ret := make([]interface{}, 0)
	for _, v := range items {
		ret = append(ret, map[string]interface{}{
			exampleExampleDatasetsIDFieldName:          v.GetId(),
			exampleExampleDatasetsNameFieldName:        v.GetName(),
			exampleExampleDatasetsDescriptionFieldName: v.GetDescription(),
			exampleExampleDatasetsCreatedAtFieldName:   v.GetCreatedAt().String(),
			exampleExampleDatasetsGuideFieldName:       v.GetGuide(),
		})
	}
	return ret
}
