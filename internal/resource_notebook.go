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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	nb "github.com/arangodb-managed/apis/notebook/v1"
)

const (
	// Notebook field names
	notebookDeploymentIdFieldName               = "deployment_id"
	notebookURLFieldName                        = "url"
	notebookNameFieldName                       = "name"
	notebookDescriptionFieldName                = "description"
	notebookIsPausedFieldName                   = "is_paused"
	notebookLastPausedAtFieldName               = "last_paused_at"
	notebookLastResumedAtFieldName              = "last_resumed_at"
	notebookCreatedByIdFieldName                = "created_by_id"
	notebookCreatedAtFieldName                  = "created_at"
	notebookModelFieldName                      = "model"
	notebookModelIdFieldName                    = "notebook_model_id"
	notebookModelDiskSizeFieldName              = "disk_size"
	notebookIsDeletedFieldName                  = "is_deleted"
	notebookDeletedAtFieldName                  = "deleted_at"
	notebookStatusFieldName                     = "status"
	notebookStatusPhaseFieldName                = "phase"
	notebookStatusMessageFieldName              = "message"
	notebookStatusLastUpdatedAtFieldName        = "last_updated_at"
	notebookStatusEndpointFieldName             = "endpoint"
	notebookStatusUsageFieldName                = "usage"
	notebookStatusUsageLastMemoryUsageFieldName = "last_memory_usage"
	notebookStatusUsageLastCpuUsageFieldName    = "last_cpu_usage"
	notebookStatusUsageLastMemoryLimitFieldName = "last_memory_limit"
	notebookStatusUsageLastCpuLimitFieldName    = "last_cpu_limit"
)

// resourceNotebook defines a Jupyter Notebook Oasis resource.
func resourceNotebook() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Notebook Resource",

		CreateContext: resourceNotebookCreate,
		ReadContext:   resourceNotebookRead,
		UpdateContext: resourceNotebookUpdate,
		DeleteContext: resourceNotebookDelete,
		Schema: map[string]*schema.Schema{
			notebookDeploymentIdFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Deployment ID field",
				Required:    true,
			},
			notebookURLFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook URL field",
				Computed:    true,
				Optional:    true,
			},
			notebookNameFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Name field",
				Required:    true,
			},
			notebookDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Description field",
				Optional:    true,
			},
			notebookIsPausedFieldName: {
				Type:        schema.TypeBool,
				Description: "Notebook Resource Notebook Is Paused field",
				Computed:    true,
			},
			notebookLastPausedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Last Paused field",
				Computed:    true,
			},
			notebookLastResumedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Last Resumed field",
				Computed:    true,
			},
			notebookCreatedByIdFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Created By Id field",
				Computed:    true,
			},
			notebookCreatedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Created At field",
				Computed:    true,
			},
			notebookModelFieldName: {
				Type:        schema.TypeList,
				Description: "Notebook Resource Notebook Model field",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						notebookModelIdFieldName: {
							Type:        schema.TypeString,
							Description: "Notebook Resource Notebook Model ID field",
							Required:    true,
						},
						notebookModelDiskSizeFieldName: {
							Type:        schema.TypeInt,
							Description: "Notebook Resource Notebook Model Disk Size field",
							Required:    true,
						},
					},
				},
			},
			notebookIsDeletedFieldName: {
				Type:        schema.TypeBool,
				Description: "Notebook Resource Notebook Is Deleted field",
				Computed:    true,
			},
			notebookDeletedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Notebook Resource Notebook Deleted At field",
				Computed:    true,
			},
			notebookStatusFieldName: {
				Type:        schema.TypeList,
				Description: "Notebook Resource Notebook Status field",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						notebookStatusPhaseFieldName: {
							Type:        schema.TypeString,
							Description: "Notebook Resource Notebook Status Phase field",
							Computed:    true,
						},
						notebookStatusMessageFieldName: {
							Type:        schema.TypeString,
							Description: "Notebook Resource Notebook Status Message field",
							Computed:    true,
						},
						notebookStatusLastUpdatedAtFieldName: {
							Type:        schema.TypeString,
							Description: "Notebook Resource Notebook Status Last Updated At field",
							Computed:    true,
						},
						notebookStatusEndpointFieldName: {
							Type:        schema.TypeString,
							Description: "Notebook Resource Notebook Status Endpoint At field",
							Computed:    true,
						},
						notebookStatusUsageFieldName: {
							Type:        schema.TypeList,
							Description: "Notebook Resource Notebook Status Usage field",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									notebookStatusUsageLastMemoryUsageFieldName: {
										Type:        schema.TypeInt,
										Description: "Notebook Resource Notebook Status Usage Last Memory Usage field",
										Computed:    true,
									},
									notebookStatusUsageLastCpuUsageFieldName: {
										Type:        schema.TypeFloat,
										Description: "Notebook Resource Notebook Status Usage Last CPU Usage field",
										Computed:    true,
									},
									notebookStatusUsageLastMemoryLimitFieldName: {
										Type:        schema.TypeInt,
										Description: "Notebook Resource Notebook Status Usage Last Memory Limit field",
										Computed:    true,
									},
									notebookStatusUsageLastCpuLimitFieldName: {
										Type:        schema.TypeFloat,
										Description: "Notebook Resource Notebook Status Usage Last CPU Limit field",
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

// resourceNotebookRead will gather information from the Terraform store for Oasis Notebook resource and display it accordingly.
func resourceNotebookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	nbc := nb.NewNotebookServiceClient(client.conn)
	notebook, err := nbc.GetNotebook(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil || notebook == nil {
		client.log.Error().Err(err).Str("notebook-id", d.Id()).Msg("Failed to find Notebook")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenNotebookResource(notebook) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// flattenNotebookResource will take a Notebook object and turn it into a flat map for terraform digestion.
func flattenNotebookResource(notebook *nb.Notebook) map[string]interface{} {
	flattened := make(map[string]interface{})
	flattened[notebookDeploymentIdFieldName] = notebook.GetDeploymentId()
	if notebook.GetUrl() != "" {
		flattened[notebookURLFieldName] = notebook.GetUrl()
	}
	flattened[notebookNameFieldName] = notebook.GetName()
	if notebook.GetDescription() != "" {
		flattened[notebookDescriptionFieldName] = notebook.GetDescription()
	}
	flattened[notebookIsPausedFieldName] = notebook.GetIsPaused()
	if notebook.GetLastPausedAt() != nil {
		flattened[notebookLastPausedAtFieldName] = notebook.GetLastPausedAt().String()
	}
	if notebook.GetLastResumedAt() != nil {
		flattened[notebookLastResumedAtFieldName] = notebook.GetLastResumedAt().String()
	}
	if notebook.GetCreatedById() != "" {
		flattened[notebookCreatedByIdFieldName] = notebook.GetCreatedById()
	}
	if notebook.GetCreatedAt() != nil {
		flattened[notebookCreatedAtFieldName] = notebook.GetCreatedAt()
	}
	flattened[notebookCreatedAtFieldName] = notebook.GetCreatedAt().String()
	if notebook.GetModel() != nil {
		flattened[notebookModelFieldName] = flattenNotebookModelSpecResource(notebook.GetModel())
	}
	flattened[notebookIsDeletedFieldName] = notebook.GetIsDeleted()
	if notebook.GetDeletedAt() != nil {
		flattened[notebookDeletedAtFieldName] = notebook.GetDeletedAt().String()
	}
	if notebook.GetStatus() != nil {
		flattened[notebookStatusFieldName] = flattenNotebookStatus(notebook.GetStatus())
	}

	return flattened
}

// flattenNotebookModelSpecResource will take a Notebook Model Spec object and turn it into a flat map for terraform digestion.
func flattenNotebookModelSpecResource(notebookModelSpec *nb.ModelSpec) []interface{} {
	model := make(map[string]interface{})
	model[notebookModelIdFieldName] = notebookModelSpec.GetNotebookModelId()
	model[notebookModelDiskSizeFieldName] = notebookModelSpec.GetDiskSize()
	return []interface{}{
		model,
	}
}

// flattenNotebookStatus will take a Notebook Status object and turn it into a flat map for terraform digestion.
func flattenNotebookStatus(notebookStatus *nb.Status) []interface{} {
	status := make(map[string]interface{})
	status[notebookStatusPhaseFieldName] = notebookStatus.GetPhase()
	status[notebookStatusMessageFieldName] = notebookStatus.GetMessage()
	status[notebookStatusLastUpdatedAtFieldName] = notebookStatus.GetLastUpdatedAt()
	status[notebookStatusEndpointFieldName] = notebookStatus.GetEndpoint()
	status[notebookStatusUsageFieldName] = flattenNotebookStatusUsage(notebookStatus.GetUsage())
	return []interface{}{
		status,
	}
}

// flattenNotebookStatusUsage will take a Notebook Status Usage object and turn it into a flat map for terraform digestion.
func flattenNotebookStatusUsage(notebookStatusUsage *nb.Status_Usage) []interface{} {
	statusUsage := make(map[string]interface{})
	statusUsage[notebookStatusUsageLastMemoryUsageFieldName] = notebookStatusUsage.GetLastMemoryUsage()
	statusUsage[notebookStatusUsageLastCpuUsageFieldName] = notebookStatusUsage.GetLastCpuUsage()
	statusUsage[notebookStatusUsageLastMemoryLimitFieldName] = notebookStatusUsage.GetLastMemoryLimit()
	statusUsage[notebookStatusUsageLastCpuLimitFieldName] = notebookStatusUsage.GetLastCpuLimit()
	return []interface{}{
		statusUsage,
	}
}

// resourceNotebookCreate will take the schema data from the Terraform config file and call the Oasis client
// to initiate a create procedure for a Jupyter Notebook. It will call helper methods to construct the necessary data
// in order to create this object.
func resourceNotebookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	nbc := nb.NewNotebookServiceClient(client.conn)
	expanded, err := expandNotebookResource(d)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	client.log.Info().Msg("Before entering Create notebook")
	result, err := nbc.CreateNotebook(client.ctxWithToken, expanded)
	client.log.Info().Msg("Creating notebook")

	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create notebook")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceNotebookRead(ctx, d, m)
}

// expandNotebookResource will take a terraform flat map schema data and turn it into an Oasis Notebook.
func expandNotebookResource(d *schema.ResourceData) (*nb.Notebook, error) {
	ret := &nb.Notebook{}
	if v, ok := d.GetOk(notebookDeploymentIdFieldName); ok {
		ret.DeploymentId = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", notebookDeploymentIdFieldName)
	}
	if v, ok := d.GetOk(notebookNameFieldName); ok {
		ret.Name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", notebookNameFieldName)
	}
	if v, ok := d.GetOk(notebookDescriptionFieldName); ok {
		ret.Description = v.(string)
	}

	if v, ok := d.GetOk(notebookModelFieldName); ok {
		ret.Model = expandNotebookModel(v.([]interface{}))
	}

	return ret, nil
}

// expandNotebookModel will take a terraform flat map schema data and decipher the monthly schedule from it.
func expandNotebookModel(s []interface{}) *nb.ModelSpec {
	ret := &nb.ModelSpec{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[notebookModelIdFieldName]; ok {
			ret.NotebookModelId = i.(string)
		}
		if i, ok := item[notebookModelDiskSizeFieldName]; ok {
			ret.DiskSize = int32(i.(int))
		}
	}
	return ret
}

// resourceNotebookDelete will delete a given Notebook resource based on the given ID
func resourceNotebookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	nbc := nb.NewNotebookServiceClient(client.conn)
	if _, err := nbc.DeleteNotebook(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("notebook-id", d.Id()).Msg("Failed to delete Notebook")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceNotebookUpdate will take a resource diff and apply changes accordingly if there are any.
func resourceNotebookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	nbc := nb.NewNotebookServiceClient(client.conn)
	notebook, err := nbc.GetNotebook(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find Notebook")
		d.SetId("")
		return diag.FromErr(err)
	}
	// Main fields
	if d.HasChange(notebookNameFieldName) {
		notebook.Name = d.Get(notebookNameFieldName).(string)
	}

	if d.HasChange(notebookDescriptionFieldName) {
		notebook.Description = d.Get(notebookDescriptionFieldName).(string)
	}

	if d.HasChange(notebookModelFieldName) {
		notebook.Model = expandNotebookModel(d.Get(notebookModelFieldName).([]interface{}))
	}

	_, err = nbc.UpdateNotebook(client.ctxWithToken, notebook)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update notebook")
		return diag.FromErr(err)
	} else {
		d.SetId(notebook.GetId())
	}
	return resourceNotebookRead(ctx, d, m)
}
