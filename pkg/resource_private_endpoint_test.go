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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	network "github.com/arangodb-managed/apis/network/v1"
)

// TestAccResourcePrivateEndpoint verifies the Oasis Private Endpoint resource is created along with the specified properties.
func TestAccResourcePrivateEndpoint(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)
	projID, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	require.NoError(t, err)

	name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyPrivateEndpoint,
		Steps: []resource.TestStep{
			{
				Config:      testPrivateEndpointConfig("", name),
				ExpectError: regexp.MustCompile("InvalidArgument desc = Project ID missing"),
			},
			{
				Config: testPrivateEndpointConfig(projID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_private_endpoint.oasis_private_endpoint_test", privateEndpointNameFieldName, name),
					resource.TestCheckResourceAttr("oasis_private_endpoint.oasis_private_endpoint_test", privateEndpointDescriptionFieldName, "Terraform generated private endpoint"),
					resource.TestCheckResourceAttr("oasis_private_endpoint.oasis_private_endpoint_test", privateEndpointDNSNamesFieldName+".#", "2"),
					resource.TestCheckResourceAttr("oasis_private_endpoint.oasis_private_endpoint_test", privateEndpointDNSNamesFieldName+".0", "test.example.com"),
					resource.TestCheckResourceAttr("oasis_private_endpoint.oasis_private_endpoint_test", privateEndpointDNSNamesFieldName+".1", "test2.example.com"),
				),
			},
		},
	})
}

// testAccCheckDestroyPrivateEndpoint verifies the Terraform oasis_private_endpoint resource cleanup.
func testAccCheckDestroyPrivateEndpoint(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_private_endpoint" {
			continue
		}
	}

	return nil
}

// TestFlattenPrivateEndpoint tests the Oasis Private Endpoint Service flattening for Terraform schema compatibility.
func TestFlattenPrivateEndpoint(t *testing.T) {
	deploymentId := acctest.RandString(10)
	privateEndpoint := &network.PrivateEndpointService{
		Name:              "test-private-endpoint",
		Description:       "test-description",
		DeploymentId:      deploymentId,
		AlternateDnsNames: []string{"test.example.com"},
		Aks: &network.PrivateEndpointService_Aks{
			ClientSubscriptionIds: []string{"sample"},
		},
	}

	expected := map[string]interface{}{
		privateEndpointNameFieldName:                    "test-private-endpoint",
		privateEndpointDescriptionFieldName:             "test-description",
		privateEndpointDeploymentFieldName:              deploymentId,
		privateEndpointDNSNamesFieldName:                []string{"test.example.com"},
		privateEndpointAzClientSubscriptionIdsFieldName: []string{"sample"},
	}

	flattened := flattenPrivateEndpointResource(privateEndpoint)
	assert.Equal(t, expected, flattened)
}

// TestExpandPrivateEndpoint tests the Oasis Private Endpoint expansion for Terraform schema compatibility.
func TestExpandPrivateEndpoint(t *testing.T) {
	deploymentId := acctest.RandString(10)
	raw := map[string]interface{}{
		privateEndpointNameFieldName:        "test-private-endpoint",
		privateEndpointDescriptionFieldName: "test-description",
		privateEndpointDeploymentFieldName:  deploymentId,
	}
	expected := &network.PrivateEndpointService{
		Name:         "test-private-endpoint",
		Description:  "test-description",
		DeploymentId: deploymentId,
	}

	s := resourcePrivateEndpoint().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedPrivateEndpoint, err := expandPrivateEndpointResource(resourceData)
	assert.NoError(t, err)

	assert.Equal(t, expected, expandedPrivateEndpoint)
}

// testPrivateEndpointConfig contains the Terraform resource definitions for testing usage
func testPrivateEndpointConfig(projID, name string) string {
	return fmt.Sprintf(`resource "oasis_deployment" "my_oneshard_deployment" {
	terms_and_conditions_accepted = "true"
	name        = "oasis_test_deployment_terraform"
	description = "Terraform Generated Deployment"
	project     = "%s"
	location {
	  region = "aks-westus2"
	}
	version {
	  db_version = "3.9.1"
	}
	configuration {
	  model      = "oneshard"
	  node_count = 3
	  maximum_node_disk_size = 20
	}
	disk_performance = "dp30"
	notification_settings {
	  email_addresses = ["test@arangodb.com"]
	}
  }

resource "oasis_private_endpoint" "oasis_private_endpoint_test" {
  name                        = "%s"
  description                 = "Terraform generated private endpoint"
  deployment                  = oasis_deployment.my_oneshard_deployment.id
  dns_names                   = ["test.example.com", "test2.example.com"]
}
`, projID, name)
}
