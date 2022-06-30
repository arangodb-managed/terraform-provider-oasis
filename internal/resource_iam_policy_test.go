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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"

	iam "github.com/arangodb-managed/apis/iam/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/assert"
)

// TestFlattenIAMPolicy tests the Oasis IAM Policy flattening for Terraform schema compatibility.
func TestFlattenIAMPolicy(t *testing.T) {
	organization := fmt.Sprintf("/Organization/%s", acctest.RandString(10))
	iamPolicy := &iam.Policy{
		ResourceUrl: organization,
		Bindings: []*iam.RoleBinding{
			{
				RoleId:   "test-role",
				MemberId: "group:300480957",
			},
		},
	}

	expected := map[string]interface{}{
		iamPolicyURLFieldName: organization,
		iamPolicyRoleBindingFieldName: []interface{}{
			map[string]interface{}{
				iamPolicyGroupFieldName: "group:300480957",
				iamPolicyRoleFieldName:  "test-role",
				iamPolicyUserFieldName:  "group:300480957",
			},
		},
	}

	flattened := flattenIAMPolicyResource(iamPolicy)
	assert.Equal(t, expected, flattened)
}

// TestExpandIAMPolicy tests the Oasis IAM Policy expansion for Terraform schema compatibility.
func TestExpandIAMPolicy(t *testing.T) {
	organization := fmt.Sprintf("/Organization/%s", acctest.RandString(10))
	raw := map[string]interface{}{
		iamPolicyURLFieldName: organization,
		iamPolicyRoleBindingFieldName: []interface{}{
			map[string]interface{}{
				iamPolicyRoleFieldName:  "test-role",
				iamPolicyGroupFieldName: "321370957",
			},
		}}

	expected := &iam.RoleBindingsRequest{
		ResourceUrl: organization,
		Bindings: []*iam.RoleBinding{
			{
				RoleId:   "test-role",
				MemberId: "group:321370957",
			},
		},
	}

	s := resourceIAMPolicy().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedIAMPolicy, err := expandToIAMPolicy(resourceData)
	assert.NoError(t, err)

	assert.Equal(t, expected, expandedIAMPolicy)
}
