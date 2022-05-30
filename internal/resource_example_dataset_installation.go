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

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	example "github.com/arangodb-managed/apis/example/v1"
)

var (
	datasetDeploymentIdFieldName       = "deployment_id"
	datasetExampleDatasetIdFieldName   = "example_dataset_id"
	datasetCreatedAtFieldName          = "created_at"
	datasetStatusFieldName             = "status"
	datasetStatusDatabaseNameFieldName = "database_name"
	datasetStatusStateFieldName        = "state"
	datasetStatusIsAvailableFieldName  = "is_available"
	datasetStatusIsFailedFieldName     = "is_failed"
)

// resourceExampleDatasetInstallation defines an Example Dataset Installation resource.
func resourceExampleDatasetInstallation() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Example Dataset Installation Resource",

		CreateContext: resourceExampleDatasetInstallationCreate,
		ReadContext:   resourceExampleDatasetInstallationRead,
		DeleteContext: resourceExampleDatasetInstallationDelete,

		Schema: map[string]*schema.Schema{
			datasetDeploymentIdFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Example Dataset Resource Deployment ID field",
				Required:    true,
				ForceNew:    true,
			},
			datasetExampleDatasetIdFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Example Dataset Resource Example Dataset ID field",
				Required:    true,
				ForceNew:    true,
			},
			datasetCreatedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Example Dataset Resource Example Dataset Created At field",
				Computed:    true,
			},
			datasetStatusFieldName: {
				Type:        schema.TypeList,
				Description: "Oasis Example Dataset Resource Example Dataset Status field",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						datasetStatusDatabaseNameFieldName: {
							Type:        schema.TypeString,
							Description: "Oasis Example Dataset Resource Example Dataset Status Database field",
							Computed:    true,
						},
						datasetStatusStateFieldName: {
							Type:        schema.TypeString,
							Description: "Oasis Example Dataset Resource Example Dataset Status State field",
							Computed:    true,
						},
						datasetStatusIsAvailableFieldName: {
							Type:        schema.TypeBool,
							Description: "Oasis Example Dataset Resource Example Dataset Status Is Available field",
							Computed:    true,
						},
						datasetStatusIsFailedFieldName: {
							Type:        schema.TypeBool,
							Description: "Oasis Example Dataset Resource Example Dataset Status Is Failed field",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceExampleDatasetInstallationCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	examplec := example.NewExampleDatasetServiceClient(client.conn)
	req := expandExampleDatasetInstallation(data)
	resp, err := examplec.CreateExampleDatasetInstallation(client.ctxWithToken, req)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create example dataset installation.")
		return diag.FromErr(err)
	}

	data.SetId(resp.GetId())
	return resourceExampleDatasetInstallationRead(ctx, data, m)
}

func expandExampleDatasetInstallation(data *schema.ResourceData) *example.ExampleDatasetInstallation {
	ret := &example.ExampleDatasetInstallation{}
	if v, ok := data.GetOk(datasetDeploymentIdFieldName); ok {
		ret.DeploymentId = v.(string)
	}
	if v, ok := data.GetOk(datasetExampleDatasetIdFieldName); ok {
		ret.ExampledatasetId = v.(string)
	}
	return ret
}

func resourceExampleDatasetInstallationRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		data.SetId("")
		return diag.FromErr(err)
	}

	examplec := example.NewExampleDatasetServiceClient(client.conn)
	response, err := examplec.GetExampleDatasetInstallation(client.ctxWithToken, &common.IDOptions{
		Id: data.Id(),
	})
	if err != nil {
		client.log.Error().Str("installation-id", data.Id()).Err(err).Msg("Failed to get example dataset installation.")
		data.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenExampleDatasetInstallation(response) {
		if err := data.Set(k, v); err != nil {
			data.SetId("")
			return diag.FromErr(err)
		}
	}
	return nil
}

// flattenExampleDatasetInstallation takes a single dataset and converts it into a terraform consumable format.
func flattenExampleDatasetInstallation(item *example.ExampleDatasetInstallation) map[string]interface{} {
	return map[string]interface{}{
		datasetDeploymentIdFieldName:     item.GetDeploymentId(),
		datasetExampleDatasetIdFieldName: item.GetExampledatasetId(),
		datasetCreatedAtFieldName:        item.GetCreatedAt().String(),
		datasetStatusFieldName:           flattenStatus(item.GetStatus()),
	}
}

// flattenStatus takes the status portion of the dataset and converts it into a terraform consumable format.
func flattenStatus(status *example.ExampleDatasetInstallation_Status) []interface{} {
	return []interface{}{
		map[string]interface{}{
			datasetStatusDatabaseNameFieldName: status.GetDatabaseName(),
			datasetStatusStateFieldName:        status.GetState(),
			datasetStatusIsFailedFieldName:     status.GetIsFailed(),
			datasetStatusIsAvailableFieldName:  status.GetIsAvailable(),
		},
	}
}

func resourceExampleDatasetInstallationDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	examplec := example.NewExampleDatasetServiceClient(client.conn)
	if _, err := examplec.DeleteExampleDatasetInstallation(client.ctxWithToken, &common.IDOptions{Id: data.Id()}); err != nil {
		client.log.Error().Err(err).Str("installation-id", data.Id()).Msg("Failed to delete installation")
		return diag.FromErr(err)
	}
	data.SetId("") // called automatically, but added to be explicit
	return nil
}
