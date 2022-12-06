//
// DISCLAIMER
//
// Copyright 2020-2022 ArangoDB GmbH, Cologne, Germany
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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
)

// TestResourceDeployment verifies the Oasis Deployment resource is created along with the specified properties.
func TestResourceDeployment(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	res := "terraform-deployment-" + acctest.RandString(10)
	name := "deployment-" + acctest.RandString(10)
	orgID, err := FetchOrganizationID()
	require.NoError(t, err)
	pid, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyDeployment,
		Steps: []resource.TestStep{
			{
				Config: testDeploymentConfig(res, name, pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_deployment."+res, deplNameFieldName, name),
					resource.TestCheckResourceAttr("oasis_deployment."+res, deplDiskPerformanceFieldName, "dp30"),
				),
			},
		},
	})
}

// TestFlattenDeploymentResource tests the Oasis Deployment flattening for Terraform schema compatibility.
func TestFlattenDeploymentResource(t *testing.T) {
	deploymentProfileTestID := acctest.RandString(10)
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "123456789",
		RegionId:    "gcp-europe-west4",
		Version:     "3.9.1",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpallowlistId: "ip-allowlist",
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
		DiskAutoSizeSettings: &data.Deployment_DiskAutoSizeSettings{
			MaximumNodeDiskSize: 40,
		},
		DiskPerformanceId:                      "dp-1",
		IsScheduledRootPasswordRotationEnabled: false,
		Locked:                                 false,
		DeploymentProfileId:                    deploymentProfileTestID,
	}
	flattened := flattenDeployment(depl)
	expected := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.9.1",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName:             "certificate-id",
				deplSecurityIpAllowlistFieldName:               "ip-allowlist",
				deplSecurityDisableFoxxAuthenticationFieldName: false,
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:               "oneshard",
				deplConfigurationNodeSizeIdFieldName:          "a8",
				deplConfigurationNodeCountFieldName:           3,
				deplConfigurationNodeDiskSizeFieldName:        32,
				deplConfigurationMaximumNodeDiskSizeFieldName: 40,
			},
		},
		deplDiskPerformanceFieldName:                      "dp-1",
		deplDisableScheduledRootPasswordRotationFieldName: true,
		deplLockedFieldName:                               false,
		deplDeploymentProfileIDFieldName:                  deploymentProfileTestID,
	}
	assert.Equal(t, expected, flattened)
}

// TestFlattenDeploymentResourceDisableFoxxAuth tests the Oasis Deployment flattening with DisableFoxxAuthentication set to true.
func TestFlattenDeploymentResourceDisableFoxxAuth(t *testing.T) {
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "123456789",
		RegionId:    "gcp-europe-west4",
		Version:     "3.9.1",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpallowlistId:             "ip-allowlist",
		DisableFoxxAuthentication: true,
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
		IsScheduledRootPasswordRotationEnabled: true,
		Locked:                                 true,
	}
	flattened := flattenDeployment(depl)
	expected := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.9.1",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName:             "certificate-id",
				deplSecurityIpAllowlistFieldName:               "ip-allowlist",
				deplSecurityDisableFoxxAuthenticationFieldName: true,
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:        "oneshard",
				deplConfigurationNodeSizeIdFieldName:   "a8",
				deplConfigurationNodeCountFieldName:    3,
				deplConfigurationNodeDiskSizeFieldName: 32,
			},
		},
		deplDiskPerformanceFieldName:                      "", // Not set
		deplDisableScheduledRootPasswordRotationFieldName: false,
		deplLockedFieldName:                               true,
	}
	assert.Equal(t, expected, flattened)
}

// TestFlattenDeploymentResourceNotificationSettings tests the Oasis Deployment flattening with notification settings.
func TestFlattenDeploymentResourceNotificationSettings(t *testing.T) {
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "123456789",
		RegionId:    "gcp-europe-west4",
		Version:     "3.9.1",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpallowlistId:             "ip-allowlist",
		DisableFoxxAuthentication: true,
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
		NotificationSettings: &data.Deployment_NotificationSettings{
			EmailAddresses: []string{"test@example.test"},
		},
		IsScheduledRootPasswordRotationEnabled: false,
		Locked:                                 true,
	}
	flattened := flattenDeployment(depl)
	expected := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.9.1",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName:             "certificate-id",
				deplSecurityIpAllowlistFieldName:               "ip-allowlist",
				deplSecurityDisableFoxxAuthenticationFieldName: true,
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:        "oneshard",
				deplConfigurationNodeSizeIdFieldName:   "a8",
				deplConfigurationNodeCountFieldName:    3,
				deplConfigurationNodeDiskSizeFieldName: 32,
			},
		},
		deplNotificationConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplNotificationConfigurationEmailAddressesFieldName: []string{"test@example.test"},
			},
		},
		deplDiskPerformanceFieldName:                      "",
		deplDisableScheduledRootPasswordRotationFieldName: true,
		deplLockedFieldName:                               true,
	}
	assert.Equal(t, expected, flattened)
}

// TestExpandingDeploymentResource tests the Oasis Deployment expansion for Terraform schema compatibility.
func TestExpandingDeploymentResource(t *testing.T) {
	deploymentProfileTestID := acctest.RandString(10)
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "123456789",
		RegionId:    "gcp-europe-west4",
		Version:     "3.9.1",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpallowlistId:             "ip-allowlist",
		DisableFoxxAuthentication: false,
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
		DiskAutoSizeSettings: &data.Deployment_DiskAutoSizeSettings{
			MaximumNodeDiskSize: 40,
		},
		DiskPerformanceId:                      "dp-2",
		IsScheduledRootPasswordRotationEnabled: true,
		Locked:                                 false,
		DeploymentProfileId:                    deploymentProfileTestID,
	}
	raw := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.9.1",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName:             "certificate-id",
				deplSecurityIpAllowlistFieldName:               "ip-allowlist",
				deplSecurityDisableFoxxAuthenticationFieldName: false,
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:               "oneshard",
				deplConfigurationNodeSizeIdFieldName:          "a8",
				deplConfigurationNodeCountFieldName:           3,
				deplConfigurationNodeDiskSizeFieldName:        32,
				deplConfigurationMaximumNodeDiskSizeFieldName: 40,
			},
		},
		deplDiskPerformanceFieldName:                      "dp-2",
		deplDisableScheduledRootPasswordRotationFieldName: false,
		deplLockedFieldName:                               false,
		deplDeploymentProfileIDFieldName:                  deploymentProfileTestID,
	}
	s := resourceDeployment().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedDepl, err := expandDeploymentResource(resourceData, "123456789")
	assert.NoError(t, err)
	assert.Equal(t, depl, expandedDepl)
}

// TestExpandingDeploymentResourceDisableFoxxAuth tests the Oasis Deployment expansion with DisableFoxxAuthentication set to true.
func TestExpandingDeploymentResourceDisableFoxxAuth(t *testing.T) {
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "123456789",
		RegionId:    "gcp-europe-west4",
		Version:     "3.9.1",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpallowlistId:             "ip-allowlist",
		DisableFoxxAuthentication: true,
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
		IsScheduledRootPasswordRotationEnabled: true,
		Locked:                                 true,
	}
	raw := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.9.1",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName:             "certificate-id",
				deplSecurityIpAllowlistFieldName:               "ip-allowlist",
				deplSecurityDisableFoxxAuthenticationFieldName: true,
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:        "oneshard",
				deplConfigurationNodeSizeIdFieldName:   "a8",
				deplConfigurationNodeCountFieldName:    3,
				deplConfigurationNodeDiskSizeFieldName: 32,
			},
		},
		deplDisableScheduledRootPasswordRotationFieldName: false,
		deplLockedFieldName: true,
	}
	s := resourceDeployment().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedDepl, err := expandDeploymentResource(resourceData, "123456789")
	assert.NoError(t, err)
	assert.Equal(t, depl, expandedDepl)
}

// TestExpandDeploymentOverrideProjectID tests the Oasis Deployment expansion when overriding project id.
func TestExpandDeploymentOverrideProjectID(t *testing.T) {
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "overrideid",
		RegionId:    "gcp-europe-west4",
		Version:     "3.9.1",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpallowlistId: "ip-allowlist",
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
		IsScheduledRootPasswordRotationEnabled: false,
		Locked:                                 true,
	}
	raw := map[string]interface{}{
		deplProjectFieldName:     "overrideid",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.9.1",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName: "certificate-id",
				deplSecurityIpAllowlistFieldName:   "ip-allowlist",
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:        "oneshard",
				deplConfigurationNodeSizeIdFieldName:   "a8",
				deplConfigurationNodeCountFieldName:    3,
				deplConfigurationNodeDiskSizeFieldName: 32,
			},
		},
		deplDisableScheduledRootPasswordRotationFieldName: true,
		deplLockedFieldName: true,
	}
	s := resourceDeployment().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedDepl, err := expandDeploymentResource(resourceData, "thisshouldbeoverriden")
	assert.NoError(t, err)
	assert.Equal(t, depl, expandedDepl)
}

// testDeploymentConfig contains the Terraform resource definitions for testing usage
func testDeploymentConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_deployment" "%s" {
	terms_and_conditions_accepted = "true"
	name        = "%s"
	description = "Terraform Generated Deployment"
	project     = "%s"
	location {
	  region = "gcp-europe-west4"
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
  }`, resource, name, project)
}

// testAccCheckDestroyDeployment verifies the Terraform oasis_deployment resource cleanup.
func testAccCheckDestroyDeployment(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	datac := data.NewDataServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_deployment" {
			continue
		}

		if _, err := datac.GetDeployment(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); !common.IsNotFound(err) {
			return fmt.Errorf("deployment still present")
		}
	}

	return nil
}
