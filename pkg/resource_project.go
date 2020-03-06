//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
)

// resourceProject defines the Project terraform resource Schema.
func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			projectNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},

			projectDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},

			projectOrganizationFieldName: {
				Type:     schema.TypeString,
				Optional: true, // overwrites plugin level settings if set
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},

			projectCreatedAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},

			projectIsDeletedFieldName: {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// resourceProjectCreate handles the creation lifecycle of the Project resource
// sets the ID of a given Project once the creation is successful. This will be stored in local terraform store.
func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	expanded, err := expandToProject(d, client.OrganizationID)
	if err != nil {
		return err
	}
	result, err := rmc.CreateProject(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create project")
		return err
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceProjectRead(d, m)
}

// expandToProject creates a project oasis structure out of a terraform schema.
func expandToProject(d *schema.ResourceData, defaultOrganization string) (*rm.Project, error) {
	var (
		name         string
		description  string
		organization string
	)
	if v, ok := d.GetOk(projectNameFieldName); ok {
		name = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", projectNameFieldName)
	}
	if v, ok := d.GetOk(projectDescriptionFieldName); ok {
		description = v.(string)
	}
	// Overwrite organization if it exists
	organization = defaultOrganization
	if v, ok := d.GetOk(projectOrganizationFieldName); ok {
		organization = v.(string)
	}
	return &rm.Project{
		Name:           name,
		Description:    description,
		OrganizationId: organization,
	}, nil
}

// flattenProjectResource flattens the project data into a map interface for easy storage.
func flattenProjectResource(project *rm.Project) map[string]interface{} {
	return map[string]interface{}{
		projectNameFieldName:         project.GetName(),
		projectDescriptionFieldName:  project.GetDescription(),
		projectOrganizationFieldName: project.GetOrganizationId(),
		projectCreatedAtFieldName:    project.GetCreatedAt().String(),
		projectIsDeletedFieldName:    project.GetIsDeleted(),
	}
}

// resourceProjectRead handles the read lifecycle of the project resource.
func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	p, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed to find project")
		d.SetId("")
		return err
	}
	if p == nil {
		client.log.Error().Str("project-id", d.Id()).Msg("Failed to find project")
		d.SetId("")
		return nil
	}

	for k, v := range flattenProjectResource(p) {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

// resourceProjectDelete will be called once the resource is destroyed.
func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	if _, err := rmc.DeleteProject(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed to delete project")
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceProjectUpdate handles the update lifecycle of the project resource.
func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	p, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed get project")
		return err
	}

	if d.HasChange(projectNameFieldName) {
		p.Name = d.Get(projectNameFieldName).(string)
	}
	if d.HasChange(projectDescriptionFieldName) {
		p.Description = d.Get(projectDescriptionFieldName).(string)
	}
	res, err := rmc.UpdateProject(client.ctxWithToken, p)
	if err != nil {
		client.log.Error().Err(err).Str("project-id", d.Id()).Msg("Failed to update project")
		return err
	}
	d.SetId(res.Id)
	return resourceProjectRead(d, m)
}
