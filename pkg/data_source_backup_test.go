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
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"

	backup "github.com/arangodb-managed/apis/backup/v1"
)

// TestFlattenBackupObject tests the Oasis Backup flattening for Terraform schema compatibility.
func TestFlattenBackupObject(t *testing.T) {
	createdAtTimeStamp, _ := types.TimestampProto(time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC))
	backup := backup.Backup{
		Id:             "test-id",
		Url:            "https://test.url",
		Name:           "test-name",
		Description:    "test-description",
		CreatedAt:      createdAtTimeStamp,
		BackupPolicyId: "test-policy-id",
		DeploymentId:   "test-dep-id",
	}

	expected := map[string]interface{}{
		backupDataSourceIdFieldName:           "test-id",
		backupDataSourceNameFieldName:         "test-name",
		backupDataSourceDescriptionFieldName:  "test-description",
		backupDataSourceURLFieldName:          "https://test.url",
		backupDataSourceCreatedAtFieldName:    "2022-01-01T01:01:01Z",
		backupDataSourceDeploymentIDFieldName: "test-dep-id",
		backupDataSourcePolicyIDFieldName:     "test-policy-id",
	}

	got := flattenBackupObject(&backup)
	assert.Equal(t, expected, got)
}
