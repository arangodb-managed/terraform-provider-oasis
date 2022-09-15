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
	"os"
	"regexp"
	"testing"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

// TestAccResourceMultiRegionBackup verifies the Oasis Multi Region Backup resource is created along with the specified properties
func TestAccResourceMultiRegionBackup(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	resourceName := "terraform-multi-region-backup-" + acctest.RandString(10)
	regionID := "gcp-us-central1"

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)
	projectID, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyMultiRegionBackup,
		Steps: []resource.TestStep{
			{
				Config:      testMultiRegionBackupConfig(projectID, resourceName, ""),
				ExpectError: regexp.MustCompile("unable to find parse field region_id"),
			},
			{
				Config: testMultiRegionBackupConfig(projectID, resourceName, regionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplNameFieldName, "oasis_multi_region_deployment"),
					resource.TestCheckResourceAttr("oasis_backup.backup", backupNameFieldName, "oasis_backup"),
					resource.TestCheckResourceAttr("oasis_multi_region_backup."+resourceName, backupRegionIDFieldName, regionID),
				),
			},
		},
	})
}

// testMultiRegionBackupConfig contains the Terraform resource definitions for testing usage
func testMultiRegionBackupConfig(project, backupResource, regionID string) string {
	return fmt.Sprintf(`
	resource "oasis_deployment" "my_oneshard_deployment" {
		terms_and_conditions_accepted = "true"
		project = "%s"
		name = "oasis_multi_region_deployment"
		location {
			region = "gcp-europe-west4"
		}
		security {
			disable_foxx_authentication = false
		}
		disk_performance = "dp30"
		configuration {
			model = "oneshard"
			node_size_id = "c4-a8"
			node_disk_size = 20
			maximum_node_disk_size = 40
		}
		notification_settings {
			email_addresses = [
			"test@arangodb.com"
			]
		}
	}

	resource "oasis_backup" "backup" {
		name = "oasis_backup"
		description = "test backup description update from terraform"
		deployment_id = oasis_deployment.my_oneshard_deployment.id
		upload = true
		auto_deleted_at = 20
	}

	resource "oasis_multi_region_backup" "%s" {
		source_backup_id = oasis_backup.backup.id
		region_id = "%s"
	}
`, project, backupResource, regionID)
}

// testAccCheckDestroyMultiRegionBackup verifies the Terraform oasis_multi_region_backup resource cleanup.
func testAccCheckDestroyMultiRegionBackup(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	backupc := backup.NewBackupServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_multi_region_backup" {
			continue
		}

		if _, err := backupc.GetBackup(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); !common.IsNotFound(err) {
			return fmt.Errorf("backup still present")
		}
	}

	return nil
}
