//
// DISCLAIMER
//
// Copyright 2020-2024 ArangoDB GmbH, Cologne, Germany
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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// Organization data source fields
	orgIdFieldName                          = "id"
	orgNameFieldName                        = "name"
	orgDescriptionFieldName                 = "description"
	orgUrlFieldName                         = "url"
	orgCreatedAtFieldName                   = "created_at"
	orgIsDeletedFieldName                   = "is_deleted"
	tierFieldName                           = "tier"
	tierIdFieldName                         = "id"
	tierNameFieldName                       = "name"
	tierHasSupportPlansFieldName            = "has_support_plans"
	tierHasBackupUploadsFieldName           = "has_backup_uploads"
	tierRequiresTermsAndConditionsFieldName = "requires_terms_and_conditions"
)

// dataSourceOasisOrganization defines an Organization datasource terraform type.
func dataSourceOasisOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Organization Data Source",

		ReadContext: dataSourceOasisOrganizationRead,

		Schema: map[string]*schema.Schema{
			orgIdFieldName: {
				Type:        schema.TypeString,
				Description: "Organization Data Source Organization ID field",
				Required:    true,
			},
			orgNameFieldName: {
				Type:        schema.TypeString,
				Description: "Organization Data Source Organization Name field",
				Optional:    true,
			},
			orgDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Organization Data Source Organization Description field",
				Optional:    true,
			},
			orgUrlFieldName: {
				Type:        schema.TypeString,
				Description: "Organization Data Source Organization URL field",
				Computed:    true,
			},
			orgCreatedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Organization Data Source Organization Created At field",
				Computed:    true,
			},
			orgIsDeletedFieldName: {
				Type:        schema.TypeBool,
				Description: "Organization Data Source Organization is deleted field",
				Computed:    true,
			},
			tierFieldName: {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Organization Data Source Organization Tier field",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tierIdFieldName: {
							Type:        schema.TypeString,
							Description: "Organization Data Source Organization Tier ID field",
							Computed:    true,
						},
						tierNameFieldName: {
							Type:        schema.TypeString,
							Description: "Organization Data Source Organization Tier Name field",
							Computed:    true,
						},
						tierHasSupportPlansFieldName: {
							Type:        schema.TypeBool,
							Description: "Organization Data Source Organization Tier Has Support Plans field",
							Computed:    true,
						},
						tierHasBackupUploadsFieldName: {
							Type:        schema.TypeBool,
							Description: "Organization Data Source Organization Tier Has Backup Uploads field",
							Computed:    true,
						},
						tierRequiresTermsAndConditionsFieldName: {
							Type:        schema.TypeBool,
							Description: "Organization Data Source Organization Tier Requires Terms and Conditions field",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceOasisOrganizationRead reloads the resource object from the terraform store.
func dataSourceOasisOrganizationRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	oid := data.Get(orgIdFieldName).(string)
	org, err := rmc.GetOrganization(client.ctxWithToken, &common.IDOptions{Id: oid})
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range flattenOrganizationObject(org) {
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	data.SetId(org.GetId())
	return nil
}

// flattenOrganizationObject creates a map from an Oasis Organization structure for the terraform schema.
func flattenOrganizationObject(org *rm.Organization) map[string]interface{} {
	ret := map[string]interface{}{
		orgIdFieldName:          org.GetId(),
		orgNameFieldName:        org.GetName(),
		orgDescriptionFieldName: org.GetDescription(),
		orgUrlFieldName:         org.GetUrl(),
		orgCreatedAtFieldName:   org.GetCreatedAt().AsTime().Format(time.RFC3339Nano),
		tierFieldName:           flattenTierObject(org.GetTier()),
	}

	return ret
}

// flattenTierObject will produce a schema set which can be interpreted by terraform for type safety.
// Contrary to first look, tier is not handled as TypeMap because in terraform 0.12, there is no support
// for different embedded types. The hash for the set is based on the definition of the schema for tier.
func flattenTierObject(tier *rm.Tier) *schema.Set {
	s := &schema.Set{
		F: schema.HashResource(dataSourceOasisOrganization().Schema[tierFieldName].Elem.(*schema.Resource)),
	}
	tierMap := map[string]interface{}{
		tierIdFieldName:                         tier.GetId(),
		tierNameFieldName:                       tier.GetName(),
		tierHasSupportPlansFieldName:            tier.GetHasSupportPlans(),
		tierHasBackupUploadsFieldName:           tier.GetHasBackupUploads(),
		tierRequiresTermsAndConditionsFieldName: tier.GetRequiresTermsAndConditions(),
	}
	s.Add(tierMap)
	return s
}
