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
	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// IAM Role field names
	iamRoleNameFieldName         = "name"
	iamRoleDescriptionFieldName  = "description"
	iamRoleOrganizationFieldName = "organization"
	iamRolePermissionsFieldName  = "permissions"
)

// resourceIAMRole defines an IAM Role Oasis resource.
func resourceIAMRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMRoleCreate,
		ReadContext:   resourceIAMRoleRead,
		UpdateContext: resourceIAMRoleUpdate,
		DeleteContext: resourceIAMRoleDelete,
		Schema: map[string]*schema.Schema{
			iamRoleNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			iamRoleDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			iamRoleOrganizationFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			iamRolePermissionsFieldName: {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// resourceIAMRoleCreate handles the creation lifecycle of the IAM Role resource and
// sets the ID of a given IAM Role once the creation is successful. This will be stored in local Terraform store.
func resourceIAMRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	iamc := iam.NewIAMServiceClient(client.conn)
	expanded, err := expandToIAMRole(d)
	if err != nil {
		return diag.FromErr(err)
	}
	result, err := iamc.CreateRole(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create IAM Role")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceIAMRoleRead(ctx, d, m)
}

// expandToIAMRole creates an IAM Role Oasis Resource structure out of a Terraform schema.
func expandToIAMRole(d *schema.ResourceData) (*iam.Role, error) {
	role := &iam.Role{}
	if v, ok := d.GetOk(iamRoleNameFieldName); ok {
		role.Name = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", iamRoleNameFieldName)
	}

	if v, ok := d.GetOk(iamRoleOrganizationFieldName); ok {
		role.OrganizationId = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", iamRoleOrganizationFieldName)
	}

	if v, ok := d.GetOk(iamRoleDescriptionFieldName); ok {
		role.Description = v.(string)
	}

	if v, ok := d.GetOk(iamRolePermissionsFieldName); ok {
		permissions, err := expandStringPermissionList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions
	}

	return role, nil
}

// expandStringPermissionList creates a string list of items from an interface slice. It also
// verifies if a given string item is empty or not. In case it's empty, an error is thrown.
func expandStringPermissionList(list []interface{}) ([]string, error) {
	permissions := make([]string, 0)
	for _, v := range list {
		if v, ok := v.(string); ok {
			if v == "" {
				return []string{}, fmt.Errorf("iam permissions cannot be empty")
			}
			permissions = append(permissions, v)
		}
	}
	return permissions, nil
}

// flattenIAMRoleResource flattens the IAM Role data into a map interface for easy storage.
func flattenIAMRoleResource(iamRole *iam.Role) map[string]interface{} {
	return map[string]interface{}{
		iamRoleNameFieldName:         iamRole.GetName(),
		iamRoleOrganizationFieldName: iamRole.GetOrganizationId(),
		iamRoleDescriptionFieldName:  iamRole.GetDescription(),
		iamRolePermissionsFieldName:  iamRole.GetPermissions(),
	}
}

// resourceIAMRoleRead handles the read lifecycle of the IAM Role resource.
func resourceIAMRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	p, err := iamc.GetRole(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil || p == nil {
		client.log.Error().Err(err).Str("iam-role-id", d.Id()).Msg("Failed to find IAM role")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenIAMRoleResource(p) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// resourceIAMRoleDelete will be called once the resource is destroyed.
func resourceIAMRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	if _, err := iamc.DeleteRole(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("iam-role-id", d.Id()).Msg("Failed to delete IAM Role")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceIAMRoleUpdate will take a resource diff and apply changes accordingly if there are any.
func resourceIAMRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	iamRole, err := iamc.GetRole(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to get IAM Role")
		d.SetId("")
		return diag.FromErr(err)
	}
	// Main fields
	if d.HasChange(iamRoleNameFieldName) {
		iamRole.Name = d.Get(iamRoleNameFieldName).(string)
	}
	if d.HasChange(iamRoleOrganizationFieldName) {
		iamRole.OrganizationId = d.Get(iamRoleOrganizationFieldName).(string)
	}
	if d.HasChange(iamRoleDescriptionFieldName) {
		iamRole.Description = d.Get(iamRoleDescriptionFieldName).(string)
	}
	if d.HasChange(iamRolePermissionsFieldName) {
		permissions, err := expandStringPermissionList(d.Get(iamRolePermissionsFieldName).([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		iamRole.Permissions = permissions
	}

	res, err := iamc.UpdateRole(client.ctxWithToken, iamRole)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update IAM Role")
		return diag.FromErr(err)
	} else {
		d.SetId(res.GetId())
	}
	return resourceIAMRoleRead(ctx, d, m)
}
