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
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// Organization Invite fields
	organizationInviteEmailFieldName        = "email"
	organizationInviteOrganizationFieldName = "organization"
)

// resourceOrganizationInvite defines the Organization Invite Terraform resource Schema.
func resourceOrganizationInvite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationInviteCreate,
		ReadContext:   resourceOrganizationInviteRead,
		DeleteContext: resourceOrganizationInviteDelete,

		Schema: map[string]*schema.Schema{
			organizationInviteEmailFieldName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			organizationInviteOrganizationFieldName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceOrganizationInviteCreate handles the creation lifecycle of the Organization Invite resource and
// sets the ID of a given Organization Invite once the creation is successful. This will be stored in local Terraform store.
func resourceOrganizationInviteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	expanded, err := expandToOrganizationInvite(d)
	if err != nil {
		return diag.FromErr(err)
	}
	result, err := rmc.CreateOrganizationInvite(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create organization invite")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceOrganizationInviteRead(ctx, d, m)
}

// expandToOrganizationInvite creates an Organization Resource Oasis structure out of a Terraform schema.
func expandToOrganizationInvite(d *schema.ResourceData) (*rm.OrganizationInvite, error) {
	organizationInvite := &rm.OrganizationInvite{}
	if v, ok := d.GetOk(organizationInviteOrganizationFieldName); ok {
		organizationInvite.OrganizationId = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", organizationInviteOrganizationFieldName)
	}

	if v, ok := d.GetOk(organizationInviteEmailFieldName); ok {
		organizationInvite.Email = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", organizationInviteEmailFieldName)
	}

	return organizationInvite, nil
}

// flattenOrganizationInviteResource flattens the Organization Invite data into a map interface for easy storage.
func flattenOrganizationInviteResource(organizationInvite *rm.OrganizationInvite) map[string]interface{} {
	return map[string]interface{}{
		organizationInviteEmailFieldName:        organizationInvite.GetEmail(),
		organizationInviteOrganizationFieldName: organizationInvite.GetOrganizationId(),
	}
}

// resourceOrganizationInviteRead handles the read lifecycle of the Organization Invite resource.
func resourceOrganizationInviteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	p, err := rmc.GetOrganizationInvite(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("organization-invite-id", d.Id()).Msg("Failed to find organization invite")
		d.SetId("")
		return diag.FromErr(err)
	}
	if p == nil {
		client.log.Error().Str("organization-invite-id", d.Id()).Msg("Failed to find organization invite")
		d.SetId("")
		return nil
	}

	for k, v := range flattenOrganizationInviteResource(p) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// resourceOrganizationInviteDelete will be called once the resource is destroyed.
func resourceOrganizationInviteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	if _, err := rmc.DeleteOrganizationInvite(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("organization-invite-id", d.Id()).Msg("Failed to delete organization invite")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}
