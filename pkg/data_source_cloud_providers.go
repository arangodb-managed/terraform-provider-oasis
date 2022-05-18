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

package pkg

import (
	"context"
	"crypto/sha256"
	"fmt"
	common "github.com/arangodb-managed/apis/common/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	platform "github.com/arangodb-managed/apis/platform/v1"
)

const (
	// provider data source fields
	providerOrganizationFieldName = "organization"
	providerProvidersFieldName    = "providers"
)

// dataSourceOasisCloudProvider defines a Cloud Provider datasource terraform type.
func dataSourceOasisCloudProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOasisCloudProviderRead,

		Schema: map[string]*schema.Schema{
			providerOrganizationFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			providerProvidersFieldName: {
				Type:        schema.TypeList,
				Description: "List of all supported cloud providers in Oasis.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

	providers := make([]string, len(providersRaw.Items))

	fmt.Println("prov and len", providers, len(providers))

	for i, provider := range providersRaw.Items {
		providers[i] = provider.Name
	}

	err = data.Set(providerProvidersFieldName, providers)
	if err != nil {
		return diag.FromErr(err)
	}

	idsum := sha256.New()
	for _, v := range providers {
		_, err := idsum.Write([]byte(v))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	id := fmt.Sprintf("%x", idsum.Sum(nil))
	data.SetId(id)

	return nil
}
