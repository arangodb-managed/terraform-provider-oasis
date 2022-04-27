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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// resourceOrganizationLock defines a Lock resource to lock an Oasis Organization.
func resourceOrganizationLock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationLockCreate,
		ReadContext:   resourceOrganizationLockRead,
		DeleteContext: resourceOrganizationLockDelete,

		Schema: map[string]*schema.Schema{
			organizationIdNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceOrganizationLockRead will gather information from the Terraform store for Oasis Organization locking resource and display it accordingly.
func resourceOrganizationLockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	for k, v := range flattenResourceLockedOrganizationResource(organization) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// resourceOrganizationLockCreate will take the schema data from the Terraform config file and call the Oasis client
// to initiate a create procedure for an Organization Lock. It will call helper methods to construct the necessary data
// in order to create this object.
func resourceOrganizationLockCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	expanded, err := expandResourceLockedOrganizationResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(expanded.Id)

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	organization, err := rmc.GetOrganization(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find Organization")
		d.SetId("")
		return diag.FromErr(err)
	}

	organization.Locked = true

	res, err := rmc.UpdateOrganization(client.ctxWithToken, organization)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update Organization")
		return diag.FromErr(err)
	} else {
		d.SetId(res.GetId())
	}
	return resourceOrganizationLockRead(ctx, d, m)
}

// resourceOrganizationLockDelete will delete a given resource.
func resourceOrganizationLockDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
