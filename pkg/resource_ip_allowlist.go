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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	security "github.com/arangodb-managed/apis/security/v1"
)

const (
	// IP Allowlist fields
	ipNameFieldName                    = "name"
	ipProjectFieldName                 = "project"
	ipDescriptionFieldName             = "description"
	ipCIDRRangeFieldName               = "cidr_ranges"
	ipIsDeletedFieldName               = "is_deleted"
	ipCreatedAtFieldName               = "created_at"
	ipRemoteInspectionAllowedFieldName = "remote_inspection_allowed"
)

// resourceIPAllowlist defines the IPAllowlist terraform resource Schema.
func resourceIPAllowlist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPAllowlistCreate,
		ReadContext:   resourceIPAllowlistRead,
		UpdateContext: resourceIPAllowlistUpdate,
		DeleteContext: resourceIPAllowlistDelete,

		Schema: map[string]*schema.Schema{
			ipNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},

			ipProjectFieldName: { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			ipDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},

			ipCIDRRangeFieldName: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			ipRemoteInspectionAllowedFieldName: {
				Type:     schema.TypeBool,
				Optional: true,
			},

			ipIsDeletedFieldName: {
				Type:     schema.TypeBool,
				Computed: true,
			},

			ipCreatedAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceIPAllowlistCreate handles the creation lifecycle of the IPAllowlist resource
// sets the ID of a given IPAllowlist once the creation is successful. This will be stored in local terraform store.
func resourceIPAllowlistCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	securityc := security.NewSecurityServiceClient(client.conn)
	expanded, err := expandToIPAllowlist(d, client.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}
	result, err := securityc.CreateIPAllowlist(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create ip allowlist")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceIPAllowlistRead(ctx, d, m)
}

// expandToIPAllowlist creates an ip allowlist oasis structure out of a terraform schema.
func expandToIPAllowlist(d *schema.ResourceData, defaultProject string) (*security.IPAllowlist, error) {
	var (
		name                    string
		description             string
		cidrRange               []string
		remoteInspectionAllowed bool
		err                     error
	)
	if v, ok := d.GetOk(ipNameFieldName); ok {
		name = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", ipNameFieldName)
	}
	if v, ok := d.GetOk(ipCIDRRangeFieldName); ok {
		cidrRange, err = expandStringList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("failed to parse field %s", ipNameFieldName)
	}
	if v, ok := d.GetOk(ipRemoteInspectionAllowedFieldName); ok {
		remoteInspectionAllowed = v.(bool)
	}
	project := defaultProject
	if v, ok := d.GetOk(ipDescriptionFieldName); ok {
		description = v.(string)
	}
	// Overwrite project if it exists
	if v, ok := d.GetOk(ipProjectFieldName); ok {
		project = v.(string)
	}

	return &security.IPAllowlist{
		Name:                    name,
		Description:             description,
		ProjectId:               project,
		CidrRanges:              cidrRange,
		RemoteInspectionAllowed: remoteInspectionAllowed,
	}, nil
}

// expandStringList creates a string list of items from an interface slice. It also
// verifies if a given string item is empty or not. In case it's empty, an error is thrown.
func expandStringList(list []interface{}) ([]string, error) {
	cidr := make([]string, 0)
	for _, v := range list {
		if v, ok := v.(string); ok {
			if v == "" {
				return []string{}, fmt.Errorf("cidr range cannot be empty")
			}
			cidr = append(cidr, v)
		}
	}
	return cidr, nil
}

// flattenIPAllowlistResource flattens the ip allowlist data into a map interface for easy storage.
func flattenIPAllowlistResource(ip *security.IPAllowlist) map[string]interface{} {
	return map[string]interface{}{
		ipNameFieldName:                    ip.GetName(),
		ipProjectFieldName:                 ip.GetProjectId(),
		ipDescriptionFieldName:             ip.GetDescription(),
		ipCIDRRangeFieldName:               ip.GetCidrRanges(),
		ipRemoteInspectionAllowedFieldName: ip.GetRemoteInspectionAllowed(),
		ipCreatedAtFieldName:               ip.GetCreatedAt().String(),
		ipIsDeletedFieldName:               ip.GetIsDeleted(),
	}
}

// resourceIPAllowlistRead handles the read lifecycle of the IPAllowlist resource.
func resourceIPAllowlistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	securityc := security.NewSecurityServiceClient(client.conn)
	ipAllowlist, err := securityc.GetIPAllowlist(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("ipallowlist-id", d.Id()).Msg("Failed to find ip allowlist")
		return diag.FromErr(err)
	}
	if ipAllowlist == nil {
		client.log.Error().Str("ipallowlist-id", d.Id()).Msg("Failed to find ip allowlist")
		d.SetId("")
		return nil
	}

	for k, v := range flattenIPAllowlistResource(ipAllowlist) {
		if _, ok := d.GetOk(k); ok {
			if err := d.Set(k, v); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

// resourceIPAllowlistDelete will be called once the resource is destroyed.
func resourceIPAllowlistDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	securityc := security.NewSecurityServiceClient(client.conn)
	if _, err := securityc.DeleteIPAllowlist(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("ipallowlist-id", d.Id()).Msg("Failed to delete ip allowlist")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceIPAllowlistUpdate handles the update lifecycle of the IPAllowlist resource.
func resourceIPAllowlistUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	securityc := security.NewSecurityServiceClient(client.conn)
	ipAllowlist, err := securityc.GetIPAllowlist(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("ipallowlist-id", d.Id()).Msg("Failed get ip allowlist")
		return diag.FromErr(err)
	}
	if ipAllowlist == nil {
		client.log.Error().Str("ipallowlist-id", d.Id()).Msg("Failed to find certificate")
		d.SetId("")
		return nil
	}

	if d.HasChange(ipNameFieldName) {
		ipAllowlist.Name = d.Get(ipNameFieldName).(string)
	}
	if d.HasChange(ipDescriptionFieldName) {
		ipAllowlist.Description = d.Get(ipDescriptionFieldName).(string)
	}
	if d.HasChange(ipCIDRRangeFieldName) {
		cidrRange, err := expandStringList(d.Get(ipCIDRRangeFieldName).([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		ipAllowlist.CidrRanges = cidrRange
	}
	if d.HasChange(ipRemoteInspectionAllowedFieldName) {
		ipAllowlist.RemoteInspectionAllowed = d.Get(ipRemoteInspectionAllowedFieldName).(bool)
	}
	res, err := securityc.UpdateIPAllowlist(client.ctxWithToken, ipAllowlist)
	if err != nil {
		client.log.Error().Err(err).Str("ipallowlist-id", d.Id()).Msg("Failed to update ip allowlist")
		return diag.FromErr(err)
	}
	d.SetId(res.Id)
	return resourceIPAllowlistRead(ctx, d, m)
}
