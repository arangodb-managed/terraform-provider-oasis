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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
)

// TestAccResourceIAMRole verifies the Oasis IAM Role resource is created along with the specified properties.
func TestAccResourceIAMRole(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)

	name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyIAMRole,
		Steps: []resource.TestStep{
			{
				Config:      testIAMRoleConfig("", name),
				ExpectError: regexp.MustCompile("failed to parse field organization"),
			},
			{
				Config: testIAMRoleConfig(orgID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_iam_role.oasis_iam_role_test", iamRoleOrganizationFieldName, orgID),
					resource.TestCheckResourceAttr("oasis_iam_role.oasis_iam_role_test", iamRoleNameFieldName, name),
					resource.TestCheckResourceAttr("oasis_iam_role.oasis_iam_role_test", iamRolePermissionsFieldName+".#", "1"),
					resource.TestCheckResourceAttr("oasis_iam_role.oasis_iam_role_test", iamRolePermissionsFieldName+".0", "backup.backup.list"),
				),
			},
		},
	})
}

// testAccCheckDestroyIAMRole verifies the Terraform oasis_iam_role resource cleanup.
func testAccCheckDestroyIAMRole(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	iamc := iam.NewIAMServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_iam_role" {
			continue
		}

		if _, err := iamc.GetRole(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); err == nil {
			return fmt.Errorf("iam role still present")
		}
	}

	return nil
}

// TestFlattenIAMRole tests the Oasis IAM Role flattening for Terraform schema compatibility.
func TestFlattenIAMRole(t *testing.T) {
	organizationId := acctest.RandString(10)
	iamRole := &iam.Role{
		Name:           "test-iam-role",
		Description:    "test-description",
		OrganizationId: organizationId,
		Permissions:    []string{"audit.auditlog.create"},
	}

	expected := map[string]interface{}{
		iamRoleNameFieldName:         "test-iam-role",
		iamRoleDescriptionFieldName:  "test-description",
		iamRoleOrganizationFieldName: organizationId,
		iamRolePermissionsFieldName:  []string{"audit.auditlog.create"},
	}

	flattened := flattenIAMRoleResource(iamRole)
	assert.Equal(t, expected, flattened)
}

// TestExpandIAMRole tests the Oasis IAM Role expansion for Terraform schema compatibility.
func TestExpandIAMRole(t *testing.T) {
	organizationId := acctest.RandString(10)
	raw := map[string]interface{}{
		iamRoleNameFieldName:         "test-iam-role",
		iamRoleDescriptionFieldName:  "test-description",
		iamRoleOrganizationFieldName: organizationId,
	}
	expected := &iam.Role{
		Name:           "test-iam-role",
		Description:    "test-description",
		OrganizationId: organizationId,
	}

	s := resourceIAMRole().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedIAMRole, err := expandToIAMRole(resourceData)
	assert.NoError(t, err)

	assert.Equal(t, expected, expandedIAMRole)
}

// testIAMRoleConfig contains the Terraform resource definitions for testing usage
func testIAMRoleConfig(orgID, name string) string {
	return fmt.Sprintf(`resource "oasis_iam_role" "oasis_iam_role_test" {
  organization        = "%s"
  name			 	  = "%s"
  description		  = "test description from Terraform Provider"
  permissions		  = ["backup.backup.list"]
}
`, orgID, name)
}
