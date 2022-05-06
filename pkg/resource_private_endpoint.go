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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	network "github.com/arangodb-managed/apis/network/v1"
)

const (
	// Private Endpoint field names
	privateEndpointNameFieldName                    = "name"
	privateEndpointDescriptionFieldName             = "description"
	privateEndpointDeploymentFieldName              = "deployment"
	privateEndpointDNSNamesFieldName                = "dns_names"
	privateEndpointAzClientSubscriptionIdsFieldName = "az_client_subscription_ids"
)

// resourcePrivateEndpoint defines a Private Endpoint Oasis resource.
func resourcePrivateEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateEndpointCreate,
		ReadContext:   resourcePrivateEndpointRead,
		UpdateContext: resourcePrivateEndpointUpdate,
		DeleteContext: resourcePrivateEndpointDelete,
		Schema: map[string]*schema.Schema{
			privateEndpointNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			privateEndpointDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			privateEndpointDeploymentFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			privateEndpointDNSNamesFieldName: {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			privateEndpointAzClientSubscriptionIdsFieldName: {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// resourcePrivateEndpointRead will gather information from the Terraform store for Private Endpoint resource and display it accordingly.
func resourcePrivateEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	nwc := network.NewNetworkServiceClient(client.conn)
	privateEndpoint, err := nwc.GetPrivateEndpointService(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil || privateEndpoint == nil {
		client.log.Error().Err(err).Str("private-endpoint-id", d.Id()).Msg("Failed to find Private Endpoint")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenPrivateEndpointResource(privateEndpoint) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// flattenPrivateEndpointResource will take a Private Endpoint object and turn it into a flat map for terraform digestion.
func flattenPrivateEndpointResource(privateEndpoint *network.PrivateEndpointService) map[string]interface{} {
	return map[string]interface{}{
		privateEndpointNameFieldName:                    privateEndpoint.GetName(),
		privateEndpointDescriptionFieldName:             privateEndpoint.GetDescription(),
		privateEndpointDeploymentFieldName:              privateEndpoint.GetDeploymentId(),
		privateEndpointDNSNamesFieldName:                privateEndpoint.GetAlternateDnsNames(),
		privateEndpointAzClientSubscriptionIdsFieldName: privateEndpoint.GetAks().GetClientSubscriptionIds(),
	}
}

// resourcePrivateEndpointCreate will take the schema data from the Terraform config file and call the Oasis client
// to initiate a create procedure for a Private Endpoint Service. It will call helper methods to construct the necessary data
// in order to create this object.
func resourcePrivateEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	nwc := network.NewNetworkServiceClient(client.conn)
	expanded, err := expandPrivateEndpointResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := nwc.CreatePrivateEndpointService(client.ctxWithToken, expanded)
	fmt.Println("created: ", result)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create private endpoint service")
		fmt.Println("error: ", result)
		return diag.FromErr(err)
	}

	fmt.Println("no error: ", result)
	if result != nil {
		d.SetId(result.Id)
	}
	return resourcePrivateEndpointRead(ctx, d, m)
}

// expandPrivateEndpointStringList creates a string list of items from an interface slice. It also
// verifies if a given string item is empty or not. In case it's empty, an error is thrown.
func expandPrivateEndpointStringList(list []interface{}) ([]string, error) {
	items := make([]string, 0)
	for _, v := range list {
		if v, ok := v.(string); ok {
			if v == "" {
				return []string{}, fmt.Errorf("list cannot be empty")
			}
			items = append(items, v)
		}
	}
	return items, nil
}

// expandPrivateEndpointResource will take a Terraform flat map schema data and turn it into an Oasis Private Endpoint.
func expandPrivateEndpointResource(d *schema.ResourceData) (*network.PrivateEndpointService, error) {
	ret := &network.PrivateEndpointService{}
	if v, ok := d.GetOk(privateEndpointNameFieldName); ok {
		ret.Name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", privateEndpointNameFieldName)
	}
	if v, ok := d.GetOk(privateEndpointDescriptionFieldName); ok {
		ret.Description = v.(string)
	}
	if v, ok := d.GetOk(privateEndpointDeploymentFieldName); ok {
		ret.DeploymentId = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", privateEndpointDeploymentFieldName)
	}
	if v, ok := d.GetOk(privateEndpointDNSNamesFieldName); ok {
		dnsNames, err := expandPrivateEndpointStringList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		ret.AlternateDnsNames = dnsNames
	}
	if v, ok := d.GetOk(privateEndpointAzClientSubscriptionIdsFieldName); ok {
		subscriptionIds, err := expandPrivateEndpointStringList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		ret.Aks.ClientSubscriptionIds = subscriptionIds
	}
	return ret, nil
}

// resourcePrivateEndpointDelete will delete the Terraform PrivateEndpoint resource
func resourcePrivateEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

// resourcePrivateEndpointUpdate will take a resource diff and apply changes accordingly if there are any.
func resourcePrivateEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	nwc := network.NewNetworkServiceClient(client.conn)
	privateEndpoint, err := nwc.GetPrivateEndpointService(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find Private Endpoint")
		d.SetId("")
		return diag.FromErr(err)
	}
	// Main fields
	if d.HasChange(privateEndpointNameFieldName) {
		privateEndpoint.Name = d.Get(privateEndpointNameFieldName).(string)
	}
	if d.HasChange(privateEndpointDescriptionFieldName) {
		privateEndpoint.Description = d.Get(privateEndpointDescriptionFieldName).(string)
	}
	if d.HasChange(privateEndpointDNSNamesFieldName) {
		dnsNames, err := expandPrivateEndpointStringList(d.Get(privateEndpointDNSNamesFieldName).([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		privateEndpoint.AlternateDnsNames = dnsNames
	}
	if d.HasChange(privateEndpointAzClientSubscriptionIdsFieldName) {
		subscriptionIds, err := expandPrivateEndpointStringList(d.Get(privateEndpointAzClientSubscriptionIdsFieldName).([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		privateEndpoint.Aks.ClientSubscriptionIds = subscriptionIds
	}

	_, err = nwc.UpdatePrivateEndpointService(client.ctxWithToken, privateEndpoint)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update Private Endpoint")
		return diag.FromErr(err)
	}
	return resourcePrivateEndpointRead(ctx, d, m)
}
