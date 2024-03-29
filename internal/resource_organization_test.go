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

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// TestAccResourceOrganization verifies the Oasis Organization resource is created along with the specified properties
func TestAccResourceOrganization(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	res := "terraform-organization-" + acctest.RandString(10)
	name := "organization-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyOrganization,
		Steps: []resource.TestStep{
			{
				Config:      testOrganizationConfig(res, ""),
				ExpectError: regexp.MustCompile("unable to find parse field name"),
			},
			{
				Config: testOrganizationLockedConfig(res, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_organization."+res, organizationLockFieldName, "true"),
				),
			},
			{
				Config: testOrganizationConfig(res, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_organization."+res, organizationNameFieldName, name),
					resource.TestCheckResourceAttr("oasis_organization."+res, organizationDescriptionFieldName, "A test Oasis organization from Terraform Provider"),
				),
			},
		},
	})
}

// testAccCheckDestroyOrganization verifies the Terraform oasis_organization resource cleanup.
func testAccCheckDestroyOrganization(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_organization" {
			continue
		}

		if _, err := rmc.GetOrganization(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); err == nil {
			return fmt.Errorf("organization still present")
		}
	}

	return nil
}

// testOrganizationConfig contains the Terraform resource definitions for testing usage
func testOrganizationConfig(res, name string) string {
	return fmt.Sprintf(`resource "oasis_organization" "%s" {
  name        = "%s"
  description = "A test Oasis organization from Terraform Provider"
}
`, res, name)
}

// testOrganizationLockedConfig contains the Terraform resource definitions for testing usage
func testOrganizationLockedConfig(res, name string) string {
	return fmt.Sprintf(`resource "oasis_organization" "%s" {
  name        = "%s"
  description = "A test Oasis organization from Terraform Provider"
  locked = true
}
`, res, name)
}

// TestFlattenOrganization tests the Oasis Organization flattening for Terraform schema compatibility.
func TestFlattenOrganization(t *testing.T) {
	organization := &rm.Organization{
		Name:        "test-organization",
		Description: "test-description",
	}

	expected := map[string]interface{}{
		organizationNameFieldName:        "test-organization",
		organizationDescriptionFieldName: "test-description",
	}
	t.Run("with resource locking disabled", func(tt *testing.T) {
		organization.Locked = false
		expected[organizationLockFieldName] = false

		flattened := flattenOrganizationResource(organization)
		assert.Equal(tt, expected, flattened)
	})

	t.Run("with resource locking enabled", func(tt *testing.T) {
		organization.Locked = true
		expected[organizationLockFieldName] = true

		flattened := flattenOrganizationResource(organization)
		assert.Equal(tt, expected, flattened)
	})

	t.Run("with auth providers enabled", func(tt *testing.T) {
		authProviderData := []interface{}{
			map[string]interface{}{
				enableGithubFieldName:           true,
				enableGoogleFieldName:           true,
				enableMicrosoftFieldName:        true,
				enableUsernamePasswordFieldName: true,
				enableSso:                       true,
			},
		}
		organization.AuthenticationProviders = expandAuthenticationProviders(authProviderData)
		expected[authenticationProvidersFieldName] = []interface{}{
			map[string]interface{}{
				enableGithubFieldName:           true,
				enableGoogleFieldName:           true,
				enableMicrosoftFieldName:        true,
				enableUsernamePasswordFieldName: true,
				enableSso:                       true,
			},
		}
		flattened := flattenOrganizationResource(organization)
		assert.Equal(tt, expected, flattened)
	})
}

// TestExpandOrganization tests the Oasis Organization expansion for Terraform schema compatibility.
func TestExpandOrganization(t *testing.T) {
	raw := map[string]interface{}{
		organizationNameFieldName:        "test-organization",
		organizationDescriptionFieldName: "test-description",
	}
	expected := &rm.Organization{
		Name:        "test-organization",
		Description: "test-description",
	}

	s := resourceDeployment().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedOrganization, err := expandOrganizationResource(resourceData)
	assert.NoError(t, err)

	assert.Equal(t, expected, expandedOrganization)
}
