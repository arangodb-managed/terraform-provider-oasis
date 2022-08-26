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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
)

// TestAccResourceMultiRegionBackup verifies the Oasis Multi Region Backup resource is created along with the specified properties
func TestAccResourceMultiRegionBackup(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	res := "terraform-multi-region-backup-" + acctest.RandString(10)
	name := "multi-region-backup-" + acctest.RandString(10)
	sourceBackupID := acctest.RandString((10))
	regionID := "gcp-europe-west-4"

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)
	pid, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyBackup,
		Steps: []resource.TestStep{
			{
				Config:      testMultiRegionBackupConfigIncomplete(pid, res, name, sourceBackupID),
				ExpectError: regexp.MustCompile("region ID missing"),
			},
			{
				Config:      testMultiRegionBackupConfig("", res, name, sourceBackupID, regionID),
				ExpectError: regexp.MustCompile("Project ID missing"),
			},
			{
				Config: testMultiRegionBackupConfig(pid, res, name, sourceBackupID, regionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplNameFieldName, "oasis_test_dep_tf"),
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplDiskPerformanceFieldName, "dp30"),
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplTAndCAcceptedFieldName, "true"),

					resource.TestCheckResourceAttr("oasis_multi_region_backup."+res, backupSourceBackupIDFieldName, sourceBackupID),
					resource.TestCheckResourceAttr("oasis_multi_region_backup."+res, backupRegionIDFieldName, regionID),
				),
			},
		},
	})
}

// testMultiRegionBackupConfig contains the Terraform resource definitions for testing usage
func testMultiRegionBackupConfig(project, res, name, sourceBackupID, regionID string) string {
	return fmt.Sprintf(`resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project = "%s" 
  name = "%s"
  location {
    region = "gcp-europe-west4"
  }
  version {
    db_version = "3.8.6"
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

resource "oasis_multi_region_backup" "%s" {
  source_backup_id = "%s"
  region_id = "%s"
}
`, project, res, name, sourceBackupID, regionID)
}

// testMultiRegionBackupConfigIncomplete contains the incomplete Terraform resource definitions for regression testing (expected failure)
func testMultiRegionBackupConfigIncomplete(project, res, name, sourceBackupID string) string {
	return fmt.Sprintf(`resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project = "%s" 
  name = "%s"
  location {
    region = "gcp-europe-west4"
  }
  version {
    db_version = "3.8.6"
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

resource "oasis_multi_region_backup" "%s" {
    source_backup_id = "%s"
    region_id = ""
}
`, project, res, name, sourceBackupID)
}

// testAccMultiRegionCheckDestroyBackup verifies the Terraform oasis_multi_region_backup resource cleanup.
func testAccMultiRegionCheckDestroyBackup(s *terraform.State) error {
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
