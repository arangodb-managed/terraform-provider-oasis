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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// Project fields
	projectNameFieldName         = "name"
	projectDescriptionFieldName  = "description"
	projectCreatedAtFieldName    = "created_at"
	projectOrganizationFieldName = "organization"
	projectIsDeletedFieldName    = "is_deleted"
	projectLockedFieldName       = "locked"
)

// resourceProject defines the Project terraform resource Schema.
func resourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Project Resource",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			projectNameFieldName: {
				Description: "Project Resource Project Name field",
				Type:        schema.TypeString,
				Required:    true,
			},

			projectDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Project Resource Project Description field",
				Optional:    true,
			},

			projectOrganizationFieldName: {
				Type:        schema.TypeString,
				Description: "Project Resource Organization ID field",
				Optional:    true, // overwrites plugin level settings if set
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},

			projectCreatedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Project Resource Project Created At field",
				Computed:    true,
			},

			projectIsDeletedFieldName: {
				Type:        schema.TypeBool,
				Description: "Project Resource Project IsDeleted field",
				Computed:    true,
			},
			projectLockedFieldName: {
				Type:        schema.TypeBool,
				Description: "Project Resource Project Locked field",
				Optional:    true,
			},
		},
	}
}

// resourceProjectCreate handles the creation lifecycle of the Project resource
// sets the ID of a given Project once the creation is successful. This will be stored in local terraform store.
func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	expanded, err := expandToProject(d, client.OrganizationID)
	if err != nil {
		return diag.FromErr(err)
	}
	result, err := rmc.CreateProject(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create project")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceProjectRead(ctx, d, m)
}

// expandToProject creates a project oasis structure out of a terraform schema.
func expandToProject(d *schema.ResourceData, defaultOrganization string) (*rm.Project, error) {
	proj := &rm.Project{}
	if v, ok := d.GetOk(projectNameFieldName); ok {
		proj.Name = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", projectNameFieldName)
	}
	if v, ok := d.GetOk(projectDescriptionFieldName); ok {
		proj.Description = v.(string)
	}
	if v, ok := d.GetOk(projectOrganizationFieldName); ok {
		proj.OrganizationId = v.(string)
	} else {
		proj.OrganizationId = defaultOrganization
	}
	if v, ok := d.GetOk(projectLockedFieldName); ok {
		proj.Locked = v.(bool)
	}
	return proj, nil
}

// flattenProjectResource flattens the project data into a map interface for easy storage.
func flattenProjectResource(project *rm.Project) map[string]interface{} {
	return map[string]interface{}{
		projectNameFieldName:         project.GetName(),
		projectDescriptionFieldName:  project.GetDescription(),
		projectOrganizationFieldName: project.GetOrganizationId(),
		projectCreatedAtFieldName:    project.GetCreatedAt().AsTime().Format(time.RFC3339Nano),
		projectIsDeletedFieldName:    project.GetIsDeleted(),
		projectLockedFieldName:       project.GetLocked(),
	}
}

// resourceProjectRead handles the read lifecycle of the project resource.
func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	p, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed to find project")
		d.SetId("")
		return diag.FromErr(err)
	}
	if p == nil {
		client.log.Error().Str("project-id", d.Id()).Msg("Failed to find project")
		d.SetId("")
		return nil
	}

	for k, v := range flattenProjectResource(p) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// resourceProjectDelete will be called once the resource is destroyed.
func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	if _, err := rmc.DeleteProject(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed to delete project")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceProjectUpdate handles the update lifecycle of the project resource.
func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	p, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed get project")
		return diag.FromErr(err)
	}

	if d.HasChange(projectNameFieldName) {
		p.Name = d.Get(projectNameFieldName).(string)
	}
	if d.HasChange(projectDescriptionFieldName) {
		p.Description = d.Get(projectDescriptionFieldName).(string)
	}
	if d.HasChange(projectLockedFieldName) {
		p.Locked = d.Get(projectLockedFieldName).(bool)
	}
	res, err := rmc.UpdateProject(client.ctxWithToken, p)
	if err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed to update project")
		return diag.FromErr(err)
	}
	d.SetId(res.Id)
	return resourceProjectRead(ctx, d, m)
}
