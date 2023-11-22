//
// DISCLAIMER
//
// Copyright 2022-2023 ArangoDB GmbH, Cologne, Germany
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
		EnablePrivateDns:  true,
		AlternateDnsNames: []string{"test.example.com"},
	}

	expected := map[string]interface{}{
		privateEndpointNameFieldName:              "test-private-endpoint",
		privateEndpointDescriptionFieldName:       "test-description",
		privateEndpointDeploymentFieldName:        deploymentId,
		prirvateEndpointEnablePrivateDNSFieldName: true,
		privateEndpointDNSNamesFieldName:          []string{"test.example.com"},
	}

	t.Run("flattening with aks field", func(tt *testing.T) {
		expectedAks := []interface{}{
			map[string]interface{}{
				privateEndpointAKSClientSubscriptionIdsFieldName: []string{"ba3f371b-a5e3-47bf-b097-dc3bb0a052a5"},
			},
		}
		expected[privateEndpointAKSFieldName] = expectedAks
		expected[privateEndpointAWSFieldName] = []interface{}{map[string]interface{}{
			privateEndpointAWSPrincipalFieldName: []interface{}{
				map[string]interface{}{},
			},
		}}
		var projects []string
		expected[privateEndpointGCPFieldName] = []interface{}{map[string]interface{}{
			privateEndpointGCPProjectsFieldName: projects,
		}}

		rawAks := &network.PrivateEndpointService_Aks{
			ClientSubscriptionIds: []string{"ba3f371b-a5e3-47bf-b097-dc3bb0a052a5"},
		}
		privateEndpoint.Aks = rawAks

		flattened := flattenPrivateEndpointResource(privateEndpoint)
		assert.Equal(tt, expected, flattened)
		privateEndpoint.Aks = nil
	})

	t.Run("flattening with aws field", func(tt *testing.T) {
		expectedAws := []interface{}{map[string]interface{}{
			privateEndpointAWSPrincipalFieldName: []interface{}{
				map[string]interface{}{
					privateEndpointAWSPrincipalAccountIdFieldName: "123123123123",
					privateEndpointAWSPrincipalUserNamesFieldName: []string{"test@arangodb.com"},
					privateEndpointAWSPrincipalRoleNamesFieldName: []string{"test"},
				},
			},
		}}
		var subIDs []string
		expected[privateEndpointAWSFieldName] = expectedAws
		expected[privateEndpointAKSFieldName] = []interface{}{
			map[string]interface{}{
				privateEndpointAKSClientSubscriptionIdsFieldName: subIDs,
			},
		}
		var projects []string
		expected[privateEndpointGCPFieldName] = []interface{}{map[string]interface{}{
			privateEndpointGCPProjectsFieldName: projects,
		}}

		rawAws := &network.PrivateEndpointService_Aws{
			AwsPrincipals: []*network.PrivateEndpointService_AwsPrincipals{
				{
					AccountId: "123123123123",
					UserNames: []string{"test@arangodb.com"},
					RoleNames: []string{"test"},
				},
			},
		}
		privateEndpoint.Aws = rawAws

		flattened := flattenPrivateEndpointResource(privateEndpoint)
		assert.Equal(tt, expected, flattened)
		privateEndpoint.Aws = nil
	})

	t.Run("flattening with gcp field", func(tt *testing.T) {
		expectedGcp := []interface{}{
			map[string]interface{}{
				privateEndpointGCPProjectsFieldName: []string{"project1"},
			},
		}
		expected[privateEndpointGCPFieldName] = expectedGcp
		expected[privateEndpointAWSFieldName] = []interface{}{map[string]interface{}{
			privateEndpointAWSPrincipalFieldName: []interface{}{
				map[string]interface{}{},
			},
		}}
		var subIDs []string
		expected[privateEndpointAKSFieldName] = []interface{}{map[string]interface{}{
			privateEndpointAKSClientSubscriptionIdsFieldName: subIDs,
		}}
		rawGcp := &network.PrivateEndpointService_Gcp{
			Projects: []string{"project1"},
		}
		privateEndpoint.Gcp = rawGcp

		flattened := flattenPrivateEndpointResource(privateEndpoint)
		assert.Equal(tt, expected, flattened)
		privateEndpoint.Gcp = nil
	})
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

	t.Run("expansion with aks field", func(tt *testing.T) {
		rawAks := []interface{}{
			map[string]interface{}{
				privateEndpointAKSClientSubscriptionIdsFieldName: []interface{}{"ba3f371b-a5e3-47bf-b097-dc3bb0a052a5"},
			},
		}
		raw[privateEndpointAKSFieldName] = rawAks

		expectedAks := &network.PrivateEndpointService_Aks{
			ClientSubscriptionIds: []string{"ba3f371b-a5e3-47bf-b097-dc3bb0a052a5"},
		}
		expected.Aks = expectedAks

		s := resourcePrivateEndpoint().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		expandedPrivateEndpoint, err := expandPrivateEndpointResource(resourceData)
		assert.NoError(t, err)

		assert.Equal(t, expected, expandedPrivateEndpoint)
	})

	t.Run("expansion with aws field", func(tt *testing.T) {
		rawAws := []interface{}{map[string]interface{}{
			privateEndpointAWSPrincipalFieldName: []interface{}{
				map[string]interface{}{
					privateEndpointAWSPrincipalAccountIdFieldName: "123123123123",
					privateEndpointAWSPrincipalUserNamesFieldName: []interface{}{"test@arangodb.com"},
					privateEndpointAWSPrincipalRoleNamesFieldName: []interface{}{"test"},
				},
			},
		}}
		raw[privateEndpointAWSFieldName] = rawAws

		expectedAws := &network.PrivateEndpointService_Aws{
			AwsPrincipals: []*network.PrivateEndpointService_AwsPrincipals{
				{
					AccountId: "123123123123",
					UserNames: []string{"test@arangodb.com"},
					RoleNames: []string{"test"},
				},
			},
		}
		expected.Aws = expectedAws

		s := resourcePrivateEndpoint().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		expandedPrivateEndpoint, err := expandPrivateEndpointResource(resourceData)
		assert.NoError(t, err)

		assert.Equal(t, expected, expandedPrivateEndpoint)
	})

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
  name                        	= "%s"
  description                 	= "Terraform generated private endpoint"
  deployment                  	= oasis_deployment.my_oneshard_deployment.id
  dns_names                   	= ["test.example.com", "test2.example.com"]
  aks {
    az_client_subscription_ids	= ["291bba3f-e0a5-47bc-a099-3bdcb2a50a05"]
  }
}
`, projID, name)
}
