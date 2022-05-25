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

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// project data source fields
	projIdFieldName          = "id"
	projNameFieldName        = "name"
	projDescriptionFieldName = "description"
	projUrlFieldName         = "url"
	projCreatedAtFieldName   = "created_at"
)

// dataSourceOasisProject defines a Project datasource terraform type.
func dataSourceOasisProject() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Project Data Source",
		ReadContext: dataSourceOasisProjectRead,

		Schema: map[string]*schema.Schema{
			projIdFieldName: {
				Type:        schema.TypeString,
				Description: "Project Data Source Project ID",
				Required:    true,
			},
			projNameFieldName: {
				Type:        schema.TypeString,
				Description: "Project Data Source Project Name",
				Optional:    true,
			},
			projDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Project Data Source Project Description",
				Optional:    true,
			},
			projUrlFieldName: {
				Type:        schema.TypeString,
				Description: "Project Data Source Project URL",
				Computed:    true,
			},
			projCreatedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Project Data Source Project Created At time",
				Computed:    true,
			},
		},
	}
}

// dataSourceOasisProjectRead reloads the resource object from the terraform store.
func dataSourceOasisProjectRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	pid := data.Get(projIdFieldName).(string)
	proj, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: pid})
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range flattenProjectObject(proj) {
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	data.SetId(proj.GetId())
	return nil
}

// flattenProjectObject creates a map from an Oasis Project for easy digestion by the terraform schema.
func flattenProjectObject(proj *rm.Project) map[string]interface{} {
	return map[string]interface{}{
		projIdFieldName:          proj.GetId(),
		projNameFieldName:        proj.GetName(),
		projDescriptionFieldName: proj.GetDescription(),
		projUrlFieldName:         proj.GetUrl(),
		projCreatedAtFieldName:   proj.GetCreatedAt().String(),
	}
}
