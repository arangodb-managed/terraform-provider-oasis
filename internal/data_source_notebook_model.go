//
// DISCLAIMER
//
// Copyright 2022 ArangoDB GmbH, Cologne, Germany
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

	nb "github.com/arangodb-managed/apis/notebook/v1"
)

const (
	// Notebook Model data source fields
	notebookModelDataSourceName                  = "notebook"
	notebookModelDataSourceDeploymentIdFieldName = "deployment_id"
	notebookModelDataSourceItemsFieldName        = "items"
	notebookModelDataSourceIdFieldName           = "id"
	notebookModelDataSourceNameFieldName         = "name"
	notebookModelDataSourceCPUFieldName          = "cpu"
	notebookModelDataSourceMemoryFieldName       = "memory"
	notebookModelDataSourceMaxDiskSizeFieldName  = "max_disk_size"
	notebookModelDataSourceMinDiskSizeFieldName  = "min_disk_size"
)

// dataSourceOasisNotebookModel defines a Notebook Model datasource terraform type.
func dataSourceOasisNotebookModel() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Notebook Model Data Source",

		ReadContext: dataSourceOasisNotebookModelRead,

		Schema: map[string]*schema.Schema{
			notebookModelDataSourceDeploymentIdFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Model Data Source Notebook Model Deployment ID field",
				Optional:    true,
			},
			notebookModelDataSourceItemsFieldName: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						notebookModelDataSourceIdFieldName: {
							Type:        schema.TypeString,
							Description: "Notebook Model Data Source Notebook Model ID field",
							Required:    true,
						},
						notebookModelDataSourceNameFieldName: {
							Type:        schema.TypeString,
							Description: "Notebook Model Data Source Notebook Model Name field",
							Optional:    true,
						},
						notebookModelDataSourceCPUFieldName: {
							Type:        schema.TypeFloat,
							Description: "Notebook Model Data Source Notebook Model Description field",
							Optional:    true,
						},
						notebookModelDataSourceMemoryFieldName: {
							Type:        schema.TypeInt,
							Description: "Notebook Model Data Source Notebook Model Memory field",
							Optional:    true,
						},
						notebookModelDataSourceMaxDiskSizeFieldName: {
							Type:        schema.TypeInt,
							Description: "Notebook Model Data Source Notebook Model Max Disk Size field",
							Optional:    true,
						},
						notebookModelDataSourceMinDiskSizeFieldName: {
							Type:        schema.TypeInt,
							Description: "Notebook Model Data Source Notebook Model Min Disk Size field",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceOasisNotebookModelRead reloads the resource object from the terraform store.
func dataSourceOasisNotebookModelRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	nbc := nb.NewNotebookServiceClient(client.conn)

	var deploymentId string
	if v, ok := data.GetOk(notebookModelDataSourceDeploymentIdFieldName); ok {
		deploymentId = v.(string)
	} else {
		return diag.Errorf("deployment id required")
	}

	response, err := nbc.ListNotebookModels(client.ctxWithToken, &nb.ListNotebookModelsRequest{
		DeploymentId: deploymentId,
	})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to get list of notebook models.")
		return diag.FromErr(err)
	}

	for k, v := range flattenNotebookModels(deploymentId, response.GetItems()) {
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	data.SetId(uniqueResourceID(notebookModelDataSourceName))
	return nil
}

// flattenNotebookModels takes result of datasets and converts them into a terraform consumable format.
func flattenNotebookModels(id string, items []*nb.NotebookModel) map[string]interface{} {
	return map[string]interface{}{
		notebookModelDataSourceDeploymentIdFieldName: id,
		notebookModelDataSourceItemsFieldName:        flattenNotebokModel(items),
	}
}

// flattenNotebokModel converts the list of datasets it into a Terraform consumable format.
func flattenNotebokModel(items []*nb.NotebookModel) []interface{} {
	ret := make([]interface{}, 0)
	for _, v := range items {
		ret = append(ret, map[string]interface{}{
			notebookModelDataSourceIdFieldName:          v.GetId(),
			notebookModelDataSourceNameFieldName:        v.GetName(),
			notebookModelDataSourceCPUFieldName:         v.GetCpu(),
			notebookModelDataSourceMemoryFieldName:      v.GetMemory(),
			notebookModelDataSourceMaxDiskSizeFieldName: v.GetMaxDiskSize(),
			notebookModelDataSourceMinDiskSizeFieldName: v.GetMinDiskSize(),
		})
	}
	return ret
}
