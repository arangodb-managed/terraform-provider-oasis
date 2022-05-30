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
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
)

const (
	// provider data source fields
	providerIdFieldName           = "id"
	providerNameFieldName         = "name"
	providerOrganizationFieldName = "organization"
	providerProvidersFieldName    = "providers"
)

// dataSourceOasisCloudProvider defines a Cloud Provider datasource terraform type.
func dataSourceOasisCloudProvider() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Cloud Providers Data Source",

		ReadContext: dataSourceOasisCloudProviderRead,

		Schema: map[string]*schema.Schema{
			providerOrganizationFieldName: {
				Type:        schema.TypeString,
				Description: "Cloud Provider Data Source Organization ID field",
				Required:    true,
			},
			providerProvidersFieldName: {
				Type:        schema.TypeList,
				Description: "List of all supported cloud providers in Oasis.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						providerIdFieldName: {
							Type:        schema.TypeString,
							Description: "Cloud Provider ID field",
							Computed:    true,
						},
						providerNameFieldName: {
							Type:        schema.TypeString,
							Description: "Cloud Provider Name field",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceOasisCloudProviderRead reloads the resource object from the terraform store.
func dataSourceOasisCloudProviderRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	platformc := platform.NewPlatformServiceClient(client.conn)
	organizationId := data.Get(providerOrganizationFieldName).(string)
	providersRaw, err := platformc.ListProviders(client.ctxWithToken, &platform.ListProvidersRequest{OrganizationId: organizationId, Options: &common.ListOptions{}})
	if err != nil {
		return diag.FromErr(err)
	}

	providers := make([]*platform.Provider, len(providersRaw.Items))

	providerItems := flattenCloudProviders(providersRaw.Items)
	err = data.Set(providerProvidersFieldName, providerItems)
	if err != nil {
		return diag.FromErr(err)
	}

	idsum := sha256.New()
	for _, v := range providers {
		_, err := idsum.Write([]byte(v.GetId()))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	id := fmt.Sprintf("%x", idsum.Sum(nil))
	data.SetId(id)

	return nil
}

// flattenCloudProviders takes result of datasets and converts them into a Terraform consumable format.
func flattenCloudProviders(providers []*platform.Provider) []interface{} {
	if providers != nil {
		flattened := make([]interface{}, len(providers))

		for i, providerItem := range providers {
			platformSubItems := make(map[string]interface{})

			platformSubItems["id"] = providerItem.GetId()
			platformSubItems["name"] = providerItem.GetName()

			flattened[i] = platformSubItems
		}

		return flattened
	}

	return make([]interface{}, 0)
}
