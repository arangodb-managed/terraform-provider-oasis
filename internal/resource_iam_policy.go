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
	iam "github.com/arangodb-managed/apis/iam/v1"
)

const (
	// IAM Policy field names
	iamPolicyURLFieldName         = "url"
	iamPolicyRoleBindingFieldName = "binding"
	iamPolicyRoleFieldName        = "role"
	iamPolicyGroupFieldName       = "group"
	iamPolicyUserFieldName        = "user"
)

// resourceIAMPolicy defines an IAM Policy resource.
func resourceIAMPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis IAM Policy Resource",

		CreateContext: resourceIAMPolicyCreate,
		ReadContext:   resourceIAMPolicyRead,
		DeleteContext: resourceIAMPolicyDelete,
		Schema: map[string]*schema.Schema{
			iamPolicyURLFieldName: {
				Type:        schema.TypeString,
				Description: "IAM Policy Resource IAM Policy URL",
				ForceNew:    true,
				Required:    true,
			},
			iamPolicyRoleBindingFieldName: {
				Type:        schema.TypeList,
				Description: "IAM Policy Resource IAM Policy Bindings",
				MinItems:    1,
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						iamPolicyRoleFieldName: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IAM Policy Resource IAM Policy Role",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						iamPolicyGroupFieldName: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IAM Policy Resource IAM Policy Group",
						},
						iamPolicyUserFieldName: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IAM Policy Resource IAM Policy User",
						},
					},
				},
			},
		},
	}
}

// resourceIAMPolicyCreate handles the creation lifecycle of the IAM Policy resource and
// sets the ID of a given IAM Policy once the creation is successful. This will be stored in local Terraform store.
func resourceIAMPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	iamc := iam.NewIAMServiceClient(client.conn)
	expanded, err := expandToIAMPolicy(d)
	if err != nil {
		return diag.FromErr(err)
	}
	result, err := iamc.AddRoleBindings(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create IAM Policy")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.GetResourceUrl())
	}
	return resourceIAMPolicyRead(ctx, d, m)
}

// expandToIAMPolicy creates IAM Policy Oasis Resource structure out of a Terraform schema.
func expandToIAMPolicy(d *schema.ResourceData) (*iam.RoleBindingsRequest, error) {
	policy := &iam.RoleBindingsRequest{}
	if v, ok := d.GetOk(iamPolicyURLFieldName); ok {
		policy.ResourceUrl = v.(string)
	} else {
		return nil, fmt.Errorf("failed to parse field %s", iamPolicyURLFieldName)
	}

	if v, ok := d.GetOk(iamPolicyRoleBindingFieldName); ok {
		bindings, err := expandIAMPolicyBindings(v.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("failed to expand IAM Policy Binding %s", iamPolicyRoleBindingFieldName)
		}
		policy.Bindings = bindings
	} else {
		return nil, fmt.Errorf("failed to parse field %s", iamPolicyRoleBindingFieldName)
	}

	return policy, nil
}

// expandIAMPolicyBindings gathers IAM Policy Binding data from the Terraform store
func expandIAMPolicyBindings(s []interface{}) ([]*iam.RoleBinding, error) {
	bindings := make([]*iam.RoleBinding, len(s))
	for i, v := range s {
		binding := &iam.RoleBinding{}
		item := v.(map[string]interface{})
		if role, ok := item[iamPolicyRoleFieldName]; ok {
			binding.RoleId = role.(string)
		}

		if group, ok := item[iamPolicyGroupFieldName]; ok {
			binding.MemberId = iam.CreateMemberIDFromGroupID(group.(string))
		} else if user, ok := item[iamPolicyUserFieldName]; ok {
			binding.MemberId = iam.CreateMemberIDFromUserID(user.(string))
		}
		bindings[i] = binding
	}
	return bindings, nil
}

// flattenIAMPolicyResource flattens the IAM Policy data into a map interface for easy storage.
func flattenIAMPolicyResource(policy *iam.Policy) map[string]interface{} {
	return map[string]interface{}{
		iamPolicyURLFieldName:         policy.GetResourceUrl(),
		iamPolicyRoleBindingFieldName: flattenIAMPolicyBindings(policy.GetBindings()),
	}
}

// flattenIAMPolicyBindings will take an IAM Policy Binding part of an IAM Policy and create a sub map for terraform schema.
func flattenIAMPolicyBindings(iamBindings []*iam.RoleBinding) []interface{} {
	var bindings = make(map[string]interface{})
	for _, binding := range iamBindings {
		bindings[iamPolicyRoleFieldName] = binding.GetRoleId()
		bindings[iamPolicyUserFieldName] = binding.GetMemberId()
		bindings[iamPolicyGroupFieldName] = binding.GetMemberId()
	}
	return []interface{}{
		bindings,
	}
}

// resourceIAMPolicyRead handles the read lifecycle of the IAM Policy resource.
func resourceIAMPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	p, err := iamc.GetPolicy(client.ctxWithToken, &common.URLOptions{Url: d.Id()})
	if err != nil || p == nil {
		client.log.Error().Err(err).Str("iam-policy-id", d.Id()).Msg("Failed to find IAM policy")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenIAMPolicyResource(p) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// resourceIAMPolicyDelete will be called once the resource is destroyed.
func resourceIAMPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("") // called automatically, but added to be explicit
	return nil
}
