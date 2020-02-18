//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gogo/protobuf/types"

	backup "github.com/arangodb-managed/apis/backup/v1"
)

func TestFlattenBackupPolicyWithHourlySchedule(t *testing.T) {
	policy := &backup.BackupPolicy{
		Name:         "test-policy",
		Description:  "test-description",
		DeploymentId: "123456",
		IsPaused:     true,
		Schedule: &backup.BackupPolicy_Schedule{
			ScheduleType: "hourly",
			HourlySchedule: &backup.BackupPolicy_HourlySchedule{
				ScheduleEveryIntervalHours: 10,
			},
		},
		Upload:            true,
		RetentionPeriod:   types.DurationProto(200 * 24 * time.Hour),
		EmailNotification: "None",
	}
	expected := map[string]interface{}{
		backupPolicyNameFieldName:         "test-policy",
		backupPolicyDescriptionFieldName:  "test-description",
		backupPolicyDeploymentIDFieldName: "123456",
		backupPolicyIsPausedFieldName:     true,
		backupPolicyScheduleFieldName: []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: "hourly",
				backupPolicyScheudleHourlyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheudleHourlyScheduleIntervalFieldName: 10,
					},
				},
			},
		},
		backupPolicyUploadFieldName:            true,
		backupPolicyRetentionPeriodFieldName:   200,
		backupPolictEmailNotificationFeidlName: "None",
	}
	flattened := flattenBackupPolicyResource(policy)
	assert.Equal(t, expected, flattened)
}
