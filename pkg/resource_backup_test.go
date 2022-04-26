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
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backup "github.com/arangodb-managed/apis/backup/v1"
)

// TestResourceBackup verifies the Oasis Backup resource is created along with the specified properties
func TestResourceBackup(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	res := "terraform-backup-" + acctest.RandString(10)
	name := "backup-" + acctest.RandString(10)

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
				Config: testBackupConfig(pid, res, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplNameFieldName, "oasis_test_dep_tf"),
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplDiskPerformanceFieldName, "dp30"),
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplTAndCAcceptedFieldName, "true"),

					resource.TestCheckResourceAttr("oasis_backup."+res, backupNameFieldName, name),
					resource.TestCheckResourceAttr("oasis_backup."+res, backupUploadFieldName, "true"),
					resource.TestCheckResourceAttr("oasis_backup."+res, backupAutoDeleteAtFieldName, "3"),
				),
			},
		},
	})
}

// testBackupConfig contains the Terraform resource definitions for testing usage
func testBackupConfig(project, res, name string) string {
	return fmt.Sprintf(`resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project = "%s" 
  name = "oasis_test_dep_tf"
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

resource "oasis_backup" "%s" {
  name = "%s"
  description = "test backup description update from terraform"
  deployment_id = oasis_deployment.my_oneshard_deployment.id
  upload = true
  auto_deleted_at = 3
}
`, project, res, name)
}

// TestExpandBackup tests the Oasis Backup expansion for Terraform schema compatibility.
func TestExpandBackup(t *testing.T) {
	raw := map[string]interface{}{
		backupNameFieldName:         "test-backup",
		backupDescriptionFieldName:  "test-description",
		backupDeploymentIDFieldName: "test-deployment",
	}
	expected := &backup.Backup{
		Name:         "test-backup",
		Description:  "test-description",
		DeploymentId: "test-deployment",
	}
	t.Run("test backup without cloud storage upload", func(tt *testing.T) {
		raw[backupUploadFieldName] = false

		s := resourceBackup().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		backup, err := expandBackupResource(resourceData)
		assert.NoError(t, err)

		expected.Upload = false
		assert.Equal(t, expected, backup)
	})
	t.Run("test backup with cloud storage upload", func(tt *testing.T) {
		raw[backupUploadFieldName] = true

		s := resourceBackup().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		backup, err := expandBackupResource(resourceData)
		assert.NoError(t, err)

		expected.Upload = true
		assert.Equal(t, expected, backup)
	})
	t.Run("test manual backup with auto delete days", func(tt *testing.T) {
		raw[backupAutoDeleteAtFieldName] = 6

		s := resourceBackup().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		backup, err := expandBackupResource(resourceData)
		assert.NoError(t, err)

		autoDeleteAt, err := types.TimestampProto(time.Now().AddDate(0, 0, raw[backupAutoDeleteAtFieldName].(int)))
		assert.NoError(t, err)

		expected.AutoDeletedAt = autoDeleteAt
		assert.Equal(t, expected.AutoDeletedAt.GetSeconds(), backup.AutoDeletedAt.GetSeconds())
	})
}

// TestFlattenBackup tests the Oasis Backup flattening for Terraform schema compatibility.
func TestFlattenBackup(t *testing.T) {
	backup := &backup.Backup{
		Name:           "test-backup",
		Description:    "test-description",
		DeploymentId:   "123456",
		BackupPolicyId: "456123",
		Url:            "test-url",
	}

	expected := map[string]interface{}{
		backupNameFieldName:         "test-backup",
		backupDescriptionFieldName:  "test-description",
		backupDeploymentIDFieldName: "123456",
		backupPolicyIDFieldName:     "456123",
		backupURLFieldName:          "test-url",
	}

	flattened := flattenBackupResource(backup)
	assert.Equal(t, expected, flattened)
}
