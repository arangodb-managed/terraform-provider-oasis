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

// TestAccResourceIAMGroup verifies the Oasis IAM Group resource is created along with the specified properties.
func TestAccResourceIAMGroup(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)

	fmt.Println(orgID)

	name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyIAMGroup,
		Steps: []resource.TestStep{
			{
				Config:      testIAMGroupConfig("", name),
				ExpectError: regexp.MustCompile("failed to parse field organization"),
			},
			{
				Config: testIAMGroupConfig(orgID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_iam_group.oasis_iam_group_test", iamGroupOrganizationFieldName, orgID),
					resource.TestCheckResourceAttr("oasis_iam_group.oasis_iam_group_test", iamGroupNameFieldName, name),
				),
			},
		},
	})
}

// testAccCheckDestroyIAMGroup verifies the Terraform oasis_iam_group resource cleanup.
func testAccCheckDestroyIAMGroup(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	iamc := iam.NewIAMServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_iam_group" {
			continue
		}

		if _, err := iamc.GetGroup(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); err == nil {
			return fmt.Errorf("iam group still present")
		}
	}

	return nil
}

// TestFlattenIAMGroup tests the Oasis IAM Group flattening for Terraform schema compatibility.
func TestFlattenIAMGroup(t *testing.T) {
	organizationId := acctest.RandString(10)
	iamGroup := &iam.Group{
		Name:           "test-iam-group",
		Description:    "test-description",
		OrganizationId: organizationId,
	}

	expected := map[string]interface{}{
		iamGroupNameFieldName:         "test-iam-group",
		iamGroupDescriptionFieldName:  "test-description",
		iamGroupOrganizationFieldName: organizationId,
	}

	flattened := flattenIAMGroupResource(iamGroup)
	assert.Equal(t, expected, flattened)
}

// TestExpandIAMGroup tests the Oasis IAM Group expansion for Terraform schema compatibility.
func TestExpandIAMGroup(t *testing.T) {
	organizationId := acctest.RandString(10)
	raw := map[string]interface{}{
		iamGroupNameFieldName:         "test-iam-group",
		iamGroupDescriptionFieldName:  "test-description",
		iamGroupOrganizationFieldName: organizationId,
	}
	expected := &iam.Group{
		Name:           "test-iam-group",
		Description:    "test-description",
		OrganizationId: organizationId,
	}

	s := resourceIAMGroup().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedOrganization, err := expandToIAMGroup(resourceData)
	assert.NoError(t, err)

	assert.Equal(t, expected, expandedOrganization)
}

// testIAMGroupConfig contains the Terraform resource definitions for testing usage
func testIAMGroupConfig(orgID, name string) string {
	return fmt.Sprintf(`resource "oasis_iam_group" "oasis_iam_group_test" {
  organization        = "%s"
  name			 	  = "%s"
  description		  = "test description from Terraform Provider"
}
`, orgID, name)
}
