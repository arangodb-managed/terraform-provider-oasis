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
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// Organization field names
	organizationNameFieldName        = "name"
	organizationDescriptionFieldName = "description"
	organizationLockFieldName        = "locked"
	authenticationProvidersFieldName = "authentication_providers"
	enableGithubFieldName            = "enable_github"
	enableGoogleFieldName            = "enable_google"
	enableUsernamePasswordFieldName  = "enable_username_password"
	enableMicrosoftFieldName         = "enable_microsoft"
)

// resourceOrganization defines an Organization Oasis resource.
func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Organization Resource",

		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		UpdateContext: resourceOrganizationUpdate,
		DeleteContext: resourceOrganizationDelete,
		Schema: map[string]*schema.Schema{
			organizationNameFieldName: {
				Type:        schema.TypeString,
				Description: "Organization Resource Organization Name field",
				Required:    true,
			},
			organizationDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Organization Resource Organization Description field",
				Optional:    true,
			},
			organizationLockFieldName: {
				Type:        schema.TypeBool,
				Description: "Organization Resource Organization Lock field",
				Optional:    true,
			},
			authenticationProvidersFieldName: {
				Type:        schema.TypeList,
				Description: "Authentication Provider field",
				Computed:    true,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						enableGithubFieldName: {
							Type:        schema.TypeBool,
							Description: "Organization Resource Enable Github Login field",
							Optional:    true,
							Default:     false,
						},
						enableGoogleFieldName: {
							Type:        schema.TypeBool,
							Description: "Organization Resource Enable Google Login field",
							Optional:    true,
							Default:     false,
						},
						enableUsernamePasswordFieldName: {
							Type:        schema.TypeBool,
							Description: "Organization Resource Enable Username Password Login field",
							Optional:    true,
							Default:     false,
						},
						enableMicrosoftFieldName: {
							Type:        schema.TypeBool,
							Description: "Organization Resource Enable Microsoft Login field",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
		},
	}
}

// resourceOrganizationRead will gather information from the Terraform store for Oasis Organization resource and display it accordingly.
func resourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	organization, err := rmc.GetOrganization(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil || organization == nil {
		client.log.Error().Err(err).Str("organization-id", d.Id()).Msg("Failed to find Organization")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenOrganizationResource(organization) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// resourceOrganizationCreate will take the schema data from the Terraform config file and call the Oasis client
// to initiate a create procedure for an Organization. It will call helper methods to construct the necessary data
// in order to create this object.
func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	expanded, err := expandOrganizationResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := rmc.CreateOrganization(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create organization")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceOrganizationRead(ctx, d, m)
}

// expandOrganizationResource will take a Terraform flat map schema data and turn it into an Oasis Organization.
func expandOrganizationResource(d *schema.ResourceData) (*rm.Organization, error) {
	ret := &rm.Organization{}
	if v, ok := d.GetOk(organizationNameFieldName); ok {
		ret.Name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", organizationNameFieldName)
	}
	if v, ok := d.GetOk(organizationDescriptionFieldName); ok {
		ret.Description = v.(string)
	}
	if v, ok := d.GetOk(organizationLockFieldName); ok {
		ret.Locked = v.(bool)
	}
	if v, ok := d.GetOk(authenticationProvidersFieldName); ok {
		ret.AuthenticationProviders = expandAuthenticationProviders(v.([]interface{}))
	}
	return ret, nil
}

// resourceOrganizationDelete will delete a given Organization resource based on the given ID
func resourceOrganizationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	if _, err := rmc.DeleteOrganization(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("organization-id", d.Id()).Msg("Failed to delete Organization")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceOrganizationUpdate will take a resource diff and apply changes accordingly if there are any.
func resourceOrganizationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	organization, err := rmc.GetOrganization(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find Organization")
		d.SetId("")
		return diag.FromErr(err)
	}
	// Main fields
	if d.HasChange(organizationNameFieldName) {
		organization.Name = d.Get(organizationNameFieldName).(string)
	}
	if d.HasChange(organizationDescriptionFieldName) {
		organization.Description = d.Get(organizationDescriptionFieldName).(string)
	}
	if d.HasChange(organizationLockFieldName) {
		organization.Locked = d.Get(organizationLockFieldName).(bool)
	}
	if v, ok := d.GetOk(authenticationProvidersFieldName); ok {
		organization.AuthenticationProviders = expandAuthenticationProviders(v.([]interface{}))
	}
	res, err := rmc.UpdateOrganization(client.ctxWithToken, organization)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update Organization")
		return diag.FromErr(err)
	} else {
		d.SetId(res.GetId())
	}
	return resourceOrganizationRead(ctx, d, m)
}

// flattenOrganizationResource will take an Organization object and turn it into a flat map for terraform digestion.
func flattenOrganizationResource(organization *rm.Organization) map[string]interface{} {
	result := map[string]interface{}{
		organizationNameFieldName:        organization.GetName(),
		organizationDescriptionFieldName: organization.GetDescription(),
		organizationLockFieldName:        organization.GetLocked(),
	}
	if organization.GetAuthenticationProviders() != nil {
		result[authenticationProvidersFieldName] = flattenAuthenticationProviders(organization.GetAuthenticationProviders())
	}
	return result
}

// flattenAuthenticationProviders will take a AuthenticationProviders Spec object and turn it into a flat map for terraform digestion.
func flattenAuthenticationProviders(p *rm.AuthenticationProviders) []interface{} {
	providers := make(map[string]interface{})
	providers[enableGithubFieldName] = p.GetEnableGithub()
	providers[enableGoogleFieldName] = p.GetEnableGoogle()
	providers[enableMicrosoftFieldName] = p.GetEnableMicrosoft()
	providers[enableUsernamePasswordFieldName] = p.GetEnableUsernamePassword()
	return []interface{}{
		providers,
	}
}

// expandAuthenticationProviders will take a terraform flat map schema data and turn it into an ArangoGraph AuthenticationProviders.
func expandAuthenticationProviders(p []interface{}) *rm.AuthenticationProviders {
	result := &rm.AuthenticationProviders{}
	for _, v := range p {
		item := v.(map[string]interface{})
		if i, ok := item[enableGithubFieldName]; ok {
			result.EnableGithub = i.(bool)
		}
		if i, ok := item[enableGoogleFieldName]; ok {
			result.EnableGoogle = i.(bool)
		}
		if i, ok := item[enableMicrosoftFieldName]; ok {
			result.EnableMicrosoft = i.(bool)
		}
		if i, ok := item[enableUsernamePasswordFieldName]; ok {
			result.EnableUsernamePassword = i.(bool)
		}
	}
	return result
}
