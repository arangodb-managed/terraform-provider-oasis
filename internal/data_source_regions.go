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

package internal

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
)

const (
	// region data source fields
	regionIdFieldName           = "id"
	regionProviderIdFieldName   = "provider_id"
	regionLocationFieldName     = "location"
	regionAvailableFieldName    = "available"
	regionOrganizationFieldName = "organization"
	regionRegionsFieldName      = "regions"
)

// dataSourceOasisRegion defines a Cloud Provider Region datasource terraform type.
func dataSourceOasisRegion() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Cloud Provider Regions Data Source",
		ReadContext: dataSourceOasisRegionRead,

		Schema: map[string]*schema.Schema{
			regionOrganizationFieldName: {
				Type:        schema.TypeString,
				Description: "Regions Data Source Organization ID field",
				Required:    true,
			},
			regionProviderIdFieldName: {
				Type:        schema.TypeString,
				Description: "Regions Data Source Provider ID field",
				Required:    true,
			},
			regionRegionsFieldName: {
				Type:        schema.TypeList,
				Description: "List of all supported regions for a Cloud Provider in Oasis.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						regionIdFieldName: {
							Type:        schema.TypeString,
							Description: "Regions Data Source Region ID field",
							Computed:    true,
						},
						regionProviderIdFieldName: {
							Type:        schema.TypeString,
							Description: "Regions Data Source Region Provider ID field",
							Computed:    true,
						},
						regionLocationFieldName: {
							Type:        schema.TypeString,
							Description: "Regions Data Source Region Location field",
							Computed:    true,
						},
						regionAvailableFieldName: {
							Type:        schema.TypeBool,
							Description: "Regions Data Source Region Available field",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceOasisRegionRead reloads the resource object from the Terraform store.
func dataSourceOasisRegionRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	platformc := platform.NewPlatformServiceClient(client.conn)
	organizationId := data.Get(regionOrganizationFieldName).(string)
	providerId := data.Get(regionProviderIdFieldName).(string)
	regionsRaw, err := platformc.ListRegions(client.ctxWithToken, &platform.ListRegionsRequest{OrganizationId: organizationId, ProviderId: providerId, Options: &common.ListOptions{}})
	if err != nil {
		return diag.FromErr(err)
	}

	regions := make([]*platform.Region, len(regionsRaw.Items))

	regionItems := flattenRegions(regionsRaw.Items)
	err = data.Set(regionRegionsFieldName, regionItems)
	if err != nil {
		return diag.FromErr(err)
	}

	idsum := sha256.New()
	for _, v := range regions {
		_, err := idsum.Write([]byte(v.GetId()))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	id := fmt.Sprintf("%x", idsum.Sum(nil))
	data.SetId(id)

	return nil
}

// flattenRegions takes result of datasets and converts them into a Terraform consumable format.
func flattenRegions(regions []*platform.Region) []interface{} {
	if regions != nil {
		flattened := make([]interface{}, len(regions))

		for i, regionItem := range regions {
			regionSubItems := make(map[string]interface{})

			regionSubItems["id"] = regionItem.GetId()
			regionSubItems["provider_id"] = regionItem.GetProviderId()
			regionSubItems["location"] = regionItem.GetLocation()
			regionSubItems["available"] = regionItem.GetAvailable()

			flattened[i] = regionSubItems
		}

		return flattened
	}

	return make([]interface{}, 0)
}
