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
	network "github.com/arangodb-managed/apis/network/v1"
)

const (
	// Private Endpoint field names
	privateEndpointNameFieldName        = "name"
	privateEndpointDescriptionFieldName = "description"
	privateEndpointDeploymentFieldName  = "deployment"
	privateEndpointDNSNamesFieldName    = "dns_names"

	// AKS field names
	privateEndpointAKSFieldName                      = "aks"
	privateEndpointAKSClientSubscriptionIdsFieldName = "az_client_subscription_ids"

	// AWS field names
	privateEndpointAWSFieldName                   = "aws"
	privateEndpointAWSPrincipalFieldName          = "principal"
	privateEndpointAWSPrincipalAccountIdFieldName = "account_id"
	privateEndpointAWSPrincipalUserNamesFieldName = "user_names"
	privateEndpointAWSPrincipalRoleNamesFieldName = "role_names"
)

// resourcePrivateEndpoint defines a Private Endpoint Oasis resource.
func resourcePrivateEndpoint() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Private Endpoint Resource",

		CreateContext: resourcePrivateEndpointCreate,
		ReadContext:   resourcePrivateEndpointRead,
		UpdateContext: resourcePrivateEndpointUpdate,
		DeleteContext: resourcePrivateEndpointDelete,
		Schema: map[string]*schema.Schema{
			privateEndpointNameFieldName: {
				Type:        schema.TypeString,
				Description: "Private Endpoint Resource Private Endpoint Name field",
				Required:    true,
			},
			privateEndpointDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Private Endpoint Resource Private Endpoint Description field",
				Optional:    true,
			},
			privateEndpointDeploymentFieldName: {
				Type:        schema.TypeString,
				Description: "Private Endpoint Resource Private Endpoint Deployment ID field",
				Required:    true,
			},
			privateEndpointDNSNamesFieldName: {
				Type:        schema.TypeList,
				Description: "Private Endpoint Resource Private Endpoint DNS Names field (list of dns names)",
				Optional:    true,
				MinItems:    0,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			privateEndpointAKSFieldName: {
				Type:        schema.TypeList,
				Description: "Private Endpoint Resource Private Endpoint AKS field",
				Optional:    true,
				MaxItems:    1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "1" && new == "0"
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						privateEndpointAKSClientSubscriptionIdsFieldName: {
							Type:        schema.TypeList,
							Description: "Private Endpoint Resource Private Endpoint AKS Subscription IDS field (list of subscription ids)",
							Optional:    true,
							MaxItems:    1,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			privateEndpointAWSFieldName: {
				Type:        schema.TypeList,
				Description: "Private Endpoint Resource Private Endpoint AWS field",
				Optional:    true,
				MaxItems:    1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "1" && new == "0"
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						privateEndpointAWSPrincipalFieldName: {
							Type:        schema.TypeList,
							Description: "Private Endpoint Resource Private Endpoint AWS Principal field",
							MinItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									privateEndpointAWSPrincipalAccountIdFieldName: {
										Type:        schema.TypeString,
										Description: "Private Endpoint Resource Private Endpoint AWS Principal Account Id field",
										Required:    true,
									},
									privateEndpointAWSPrincipalUserNamesFieldName: {
										Type:        schema.TypeList,
										Description: "Private Endpoint Resource Private Endpoint AWS Principal User Names field",
										Optional:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									privateEndpointAWSPrincipalRoleNamesFieldName: {
										Type:        schema.TypeList,
										Description: "Private Endpoint Resource Private Endpoint AWS Principal Role Names field",
										Optional:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
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
		privateEndpointNameFieldName:        privateEndpoint.GetName(),
		privateEndpointDescriptionFieldName: privateEndpoint.GetDescription(),
		privateEndpointDeploymentFieldName:  privateEndpoint.GetDeploymentId(),
		privateEndpointDNSNamesFieldName:    privateEndpoint.GetAlternateDnsNames(),
		privateEndpointAKSFieldName:         flattenAKSResource(privateEndpoint.GetAks()),
		privateEndpointAWSFieldName:         flattenAWSResource(privateEndpoint.GetAws()),
	}
}

// flattenAKSResource will take an AKS Resource part of a Private Endpoint and create a sub map for terraform schema.
func flattenAKSResource(privateEndpointAKS *network.PrivateEndpointService_Aks) []interface{} {
	return []interface{}{
		map[string]interface{}{
			privateEndpointAKSClientSubscriptionIdsFieldName: privateEndpointAKS.GetClientSubscriptionIds(),
		},
	}
}

// flattenAWSResource will take an AWS Resource part of a Private Endpoint and create a sub map for terraform schema.
func flattenAWSResource(privateEndpointAWS *network.PrivateEndpointService_Aws) []interface{} {
	return []interface{}{
		map[string]interface{}{
			privateEndpointAWSPrincipalFieldName: flattenAWSPrincipals(privateEndpointAWS.GetAwsPrincipals()),
		},
	}
}

// flattenAWSPrincipals will take an AWS Principal Resource part of a Private Endpoint and create a sub map for terraform schema.
func flattenAWSPrincipals(privateEndpointAWSPrincipals []*network.PrivateEndpointService_AwsPrincipals) []interface{} {
	var principals = make(map[string]interface{})
	for _, principal := range privateEndpointAWSPrincipals {
		principals[privateEndpointAWSPrincipalAccountIdFieldName] = principal.GetAccountId()
		principals[privateEndpointAWSPrincipalRoleNamesFieldName] = principal.GetRoleNames()
		principals[privateEndpointAWSPrincipalUserNamesFieldName] = principal.GetUserNames()
	}
	return []interface{}{
		principals,
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
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create private endpoint service")
		return diag.FromErr(err)
	}

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
	if v, ok := d.GetOk(privateEndpointAKSFieldName); ok {
		subscriptionIds, err := expandAKSResource(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		ret.Aks = subscriptionIds
	}
	if v, ok := d.GetOk(privateEndpointAWSFieldName); ok {
		awsResource, err := expandAWSResource(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		ret.Aws = awsResource
	}
	return ret, nil
}

// expandAKSResource gathers AKS Resource data from the terraform store
func expandAKSResource(s []interface{}) (aksResource *network.PrivateEndpointService_Aks, err error) {
	for _, v := range s {
		item := v.(map[string]interface{})
		if subscriptionIds, ok := item[privateEndpointAKSClientSubscriptionIdsFieldName]; ok {
			subscriptionIds, ok := subscriptionIds.([]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to parse field %s", privateEndpointAKSClientSubscriptionIdsFieldName)
			}
			if aksResource == nil {
				aksResource = &network.PrivateEndpointService_Aks{}
			}
			for _, addr := range subscriptionIds {
				aksResource.ClientSubscriptionIds = append(aksResource.ClientSubscriptionIds, addr.(string))
			}
		}
	}
	return
}

// expandAWSResource gathers AWS Resource data from the Terraform store
func expandAWSResource(s []interface{}) (*network.PrivateEndpointService_Aws, error) {
	awsResource := &network.PrivateEndpointService_Aws{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[privateEndpointAWSPrincipalFieldName]; ok {
			awsPrincipals, err := expandAWSPrincipal(i.([]interface{}))
			if err != nil {
				return nil, err
			}
			awsResource.AwsPrincipals = awsPrincipals
		}
	}
	return awsResource, nil
}

// expandAWSPrincipal gathers AWS Resource Principal data from the Terraform store
func expandAWSPrincipal(s []interface{}) ([]*network.PrivateEndpointService_AwsPrincipals, error) {
	principals := make([]*network.PrivateEndpointService_AwsPrincipals, len(s))
	for i, v := range s {
		principal := &network.PrivateEndpointService_AwsPrincipals{}
		item := v.(map[string]interface{})
		if accountId, ok := item[privateEndpointAWSPrincipalAccountIdFieldName]; ok {
			principal.AccountId = accountId.(string)
		}
		if roleNames, ok := item[privateEndpointAWSPrincipalRoleNamesFieldName]; ok {
			roles, ok := roleNames.([]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to parse field %s", privateEndpointAWSPrincipalRoleNamesFieldName)
			}
			for _, addr := range roles {
				principal.RoleNames = append(principal.RoleNames, addr.(string))
			}
		}
		if userNames, ok := item[privateEndpointAWSPrincipalUserNamesFieldName]; ok {
			users, ok := userNames.([]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to parse field %s", privateEndpointAWSPrincipalUserNamesFieldName)
			}
			for _, addr := range users {
				principal.UserNames = append(principal.UserNames, addr.(string))
			}
		}
		principals[i] = principal
	}
	return principals, nil
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
	if d.HasChange(privateEndpointAKSFieldName) {
		aksResource, err := expandAKSResource(d.Get(privateEndpointAKSFieldName).([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		privateEndpoint.Aks = aksResource
	}
	if d.HasChange(privateEndpointAWSFieldName) {
		awsResource, err := expandAWSResource(d.Get(privateEndpointAWSFieldName).([]interface{}))
		if err != nil {
			diag.FromErr(err)
		}
		privateEndpoint.Aws = awsResource
	}

	_, err = nwc.UpdatePrivateEndpointService(client.ctxWithToken, privateEndpoint)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update Private Endpoint")
		return diag.FromErr(err)
	}
	return resourcePrivateEndpointRead(ctx, d, m)
}
