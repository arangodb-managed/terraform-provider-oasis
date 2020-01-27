//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// project data source fields
	id          = "id"
	name        = "name"
	description = "description"
	url         = "url"
	createdAt   = "created_at"
)

// dataSourceOasisProject defines a Project datasource terraform type.
func dataSourceOasisProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOasisProjectRead,

		Schema: map[string]*schema.Schema{
			id: {
				Type:     schema.TypeString,
				Required: true,
			},
			name: {
				Type:     schema.TypeString,
				Optional: true,
			},
			description: {
				Type:     schema.TypeString,
				Optional: true,
			},
			url: {
				Type:     schema.TypeString,
				Optional: true,
			},
			createdAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceOasisProjectRead reloads the resource object from the terraform store.
func dataSourceOasisProjectRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	pid := data.Get(id).(string)
	proj, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: pid})
	if err != nil {
		return err
	}

	for k, v := range flattenProjectObject(proj) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(proj.GetId())
	return nil
}

// flattenProjectObject creates a map from an Oasis Project for easy digestion by the terraform schema.
func flattenProjectObject(proj *rm.Project) map[string]interface{} {
	return map[string]interface{}{
		id:          proj.GetId(),
		name:        proj.GetName(),
		description: proj.GetDescription(),
		url:         proj.GetUrl(),
		createdAt:   proj.GetCreatedAt().String(),
	}
}
