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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	iam "github.com/arangodb-managed/apis/iam/v1"
)

const (
	// IAM Group field names
	iamGroupNameFieldName         = "name"
	iamGroupOrganizationFieldName = "organization"
	iamGroupDescriptionFieldName  = "description"
)

// resourceIAMGroup defines an IAM Group Oasis resource.
func resourceIAMGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMGroupCreate,
		ReadContext:   resourceIAMGroupRead,
		UpdateContext: resourceIAMGroupUpdate,
		DeleteContext: resourceIAMGroupDelete,
		Schema: map[string]*schema.Schema{
			iamGroupNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			iamGroupOrganizationFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			iamGroupDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// resourceIAMGroupCreate handles the creation lifecycle of the IAM Group resource and
// sets the ID of a given IAM Group once the creation is successful. This will be stored in local Terraform store.
func resourceIAMGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	iamc := iam.NewIAMServiceClient(client.conn)
	expanded, err := expandToIAMGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}
	result, err := iamc.CreateGroup(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create IAM Group")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceIAMGroupRead(ctx, d, m)
}

// expandToIAMGroup creates an IAM Group Oasis Resource structure out of a Terraform schema.
func expandToIAMGroup(d *schema.ResourceData) (*iam.Group, error) {
	group := &iam.Group{}
	if v, ok := d.GetOk(iamGroupNameFieldName); ok {
		group.Name = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", iamGroupNameFieldName)
	}

	if v, ok := d.GetOk(iamGroupOrganizationFieldName); ok {
		group.OrganizationId = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", iamGroupOrganizationFieldName)
	}

	if v, ok := d.GetOk(iamGroupDescriptionFieldName); ok {
		group.Description = v.(string)
	}

	return group, nil
}

// flattenIAMGroupResource flattens the IAM Group data into a map interface for easy storage.
func flattenIAMGroupResource(iamGroup *iam.Group) map[string]interface{} {
	return map[string]interface{}{
		iamGroupNameFieldName:         iamGroup.GetName(),
		iamGroupOrganizationFieldName: iamGroup.GetOrganizationId(),
		iamGroupDescriptionFieldName:  iamGroup.GetDescription(),
	}
}

// resourceIAMGroupRead handles the read lifecycle of the IAM Group resource.
func resourceIAMGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	p, err := iamc.GetGroup(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil || p == nil {
		client.log.Error().Err(err).Str("iam-group-id", d.Id()).Msg("Failed to find IAM group")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenIAMGroupResource(p) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// resourceIAMGroupDelete will be called once the resource is destroyed.
func resourceIAMGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	if _, err := iamc.DeleteGroup(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("iam-group-id", d.Id()).Msg("Failed to delete IAM Group")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceIAMGroupUpdate will take a resource diff and apply changes accordingly if there are any.
func resourceIAMGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	iamGroup, err := iamc.GetGroup(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to get IAM Group")
		d.SetId("")
		return diag.FromErr(err)
	}
	// Main fields
	if d.HasChange(iamGroupNameFieldName) {
		iamGroup.Name = d.Get(iamGroupNameFieldName).(string)
	}
	if d.HasChange(iamGroupOrganizationFieldName) {
		iamGroup.OrganizationId = d.Get(iamGroupOrganizationFieldName).(string)
	}
	if d.HasChange(iamGroupDescriptionFieldName) {
		iamGroup.Description = d.Get(iamGroupDescriptionFieldName).(string)
	}

	res, err := iamc.UpdateGroup(client.ctxWithToken, iamGroup)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update IAM Group")
		return diag.FromErr(err)
	} else {
		d.SetId(res.GetId())
	}
	return resourceIAMGroupRead(ctx, d, m)
}
