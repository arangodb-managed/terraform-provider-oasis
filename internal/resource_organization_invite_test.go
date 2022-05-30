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

package internal

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
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// TestAccResourceOrganizationInvite verifies the Oasis Organization Invite resource is created along with the specified properties.
func TestAccResourceOrganizationInvite(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)

	username := acctest.RandString(7)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyOrganizationInvoice,
		Steps: []resource.TestStep{
			{
				Config:      testOrganizationInviteConfig("", username),
				ExpectError: regexp.MustCompile("failed to parse field organization"),
			},
			{
				Config: testOrganizationInviteConfig(orgID, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_organization_invite.oasis_organization_invite_test", organizationInviteEmailFieldName, username+"@arangodb.com"),
					resource.TestCheckResourceAttr("oasis_organization_invite.oasis_organization_invite_test", organizationInviteOrganizationFieldName, orgID),
				),
			},
		},
	})
}

// testAccCheckDestroyOrganization verifies the Terraform oasis_organization_invite resource cleanup.
func testAccCheckDestroyOrganizationInvoice(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_organization_invite" {
			continue
		}

		if _, err := rmc.GetOrganizationInvite(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); err == nil {
			return fmt.Errorf("organization invite still present")
		}
	}

	return nil
}

// TestFlattenOrganizationInvite tests the Oasis Organization Invite flattening for Terraform schema compatibility.
func TestFlattenOrganizationInvite(t *testing.T) {
	organizationId := acctest.RandString(10)
	testUsername := acctest.RandString(7)
	organizationInvite := &rm.OrganizationInvite{
		OrganizationId: organizationId,
		Email:          testUsername + "@arangodb.com",
	}

	expected := map[string]interface{}{
		organizationInviteOrganizationFieldName: organizationId,
		organizationInviteEmailFieldName:        testUsername + "@arangodb.com",
	}
	flattened := flattenOrganizationInviteResource(organizationInvite)
	assert.Equal(t, expected, flattened)
}

// TestExpandOrganizationInvite tests the Oasis Organization Invite expansion for Terraform schema compatibility.
func TestExpandOrganizationInvite(t *testing.T) {
	organizationId := acctest.RandString(10)
	testUsername := acctest.RandString(7)
	raw := map[string]interface{}{
		organizationInviteOrganizationFieldName: organizationId,
		organizationInviteEmailFieldName:        testUsername + "@arangodb.com",
	}
	s := resourceOrganizationInvite().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	orgInvite, err := expandToOrganizationInvite(data)
	assert.NoError(t, err)

	assert.Equal(t, raw[organizationInviteEmailFieldName], orgInvite.GetEmail())
}

// testOrganizationInviteConfig contains the Terraform resource definitions for testing usage
func testOrganizationInviteConfig(orgID, username string) string {
	return fmt.Sprintf(`resource "oasis_organization_invite" "oasis_organization_invite_test" {
  organization        = "%s"
  email			 	  = "%s@arangodb.com"
}
`, orgID, username)
}
