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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestAccOasisCloudProviderBasic(t *testing.T) {
	rxPosNum := regexp.MustCompile("^[1-9][0-9]*$")

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOasisCloudProviderConfigBasic(orgID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.oasis_cloud_provider.test_oasis_cloud_providers", "providers.#", rxPosNum),
					resource.TestCheckResourceAttrSet("data.oasis_cloud_provider.test_oasis_cloud_providers", "providers.0"),
				),
			},
		},
	})
}

func testAccOasisCloudProviderConfigBasic(organizationId string) string {
	fmt.Println("org id", organizationId)
	return fmt.Sprintf(`
data "oasis_cloud_provider" "test_oasis_cloud_providers" {
	organization = "%s"
}
`, organizationId)
}
