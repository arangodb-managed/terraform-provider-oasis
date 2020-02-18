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

func TestFlattenBackupPolicy(t *testing.T) {

	policy := &backup.BackupPolicy{
		Name:              "test-policy",
		Description:       "test-description",
		DeploymentId:      "123456",
		IsPaused:          true,
		Upload:            true,
		RetentionPeriod:   types.DurationProto(200 * 24 * time.Hour),
		EmailNotification: "None",
	}

	expected := map[string]interface{}{
		backupPolicyNameFieldName:              "test-policy",
		backupPolicyDescriptionFieldName:       "test-description",
		backupPolicyDeploymentIDFieldName:      "123456",
		backupPolicyIsPausedFieldName:          true,
		backupPolicyUploadFieldName:            true,
		backupPolicyRetentionPeriodFieldName:   200,
		backupPolictEmailNotificationFeidlName: "None",
	}

	t.Run("with hourly schedule", func(tt *testing.T) {
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: "hourly",
			HourlySchedule: &backup.BackupPolicy_HourlySchedule{
				ScheduleEveryIntervalHours: 10,
			},
		}
		policy.Schedule = schedule
		expectedSchedule := []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: "hourly",
				backupPolicyScheudleHourlyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheudleHourlyScheduleIntervalFieldName: 10,
					},
				},
			},
		}
		expected[backupPolicyScheduleFieldName] = expectedSchedule
		flattened := flattenBackupPolicyResource(policy)
		assert.Equal(t, expected, flattened)
	})
	t.Run("with daily schedule", func(tt *testing.T) {
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: "daily",
			DailySchedule: &backup.BackupPolicy_DailySchedule{
				Monday:    true,
				Tuesday:   false,
				Wednesday: false,
				Thursday:  true,
				Friday:    false,
				Saturday:  false,
				Sunday:    false,
				ScheduleAt: &backup.TimeOfDay{
					Hours:    10,
					Minutes:  10,
					TimeZone: "UTC",
				},
			},
		}
		policy.Schedule = schedule
		expectedSchedule := []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: "daily",
				backupPolicyScheudleDailyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheudleDailyScheduleMondayFieldName:    true,
						backupPolicyScheudleDailyScheduleTuesdayFieldName:   false,
						backupPolicyScheudleDailyScheduleWednesdayFieldName: false,
						backupPolicyScheudleDailyScheduleThursdayFieldName:  true,
						backupPolicyScheudleDailyScheduleFridayFieldName:    false,
						backupPolicyScheudleDailyScheduleSaturdayFieldName:  false,
						backupPolicyScheudleDailyScheduleSundayFieldName:    false,
						backupPolicyTimeOfDayScheduleAtFieldName: []interface{}{
							map[string]interface{}{
								backupPolicyTimeOfDayHoursFieldName:    10,
								backupPolicyTimeOfDayMinutesFieldName:  10,
								backupPolicyTimeOfDayTimeZoneFieldName: "UTC",
							},
						},
					},
				},
			},
		}
		expected[backupPolicyScheduleFieldName] = expectedSchedule
		flattened := flattenBackupPolicyResource(policy)
		assert.Equal(t, expected, flattened)
	})
	t.Run("with monthly schedule", func(tt *testing.T) {
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: "monthly",
			MonthlySchedule: &backup.BackupPolicy_MonthlySchedule{
				DayOfMonth: 30,
				ScheduleAt: &backup.TimeOfDay{
					Hours:    10,
					Minutes:  10,
					TimeZone: "UTC",
				},
			},
		}
		policy.Schedule = schedule
		expectedSchedule := []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: "monthly",
				backupPolicyScheudleMonthlyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheudleMonthlyScheduleDayOfMonthScheduleFieldName: 30,
						backupPolicyTimeOfDayScheduleAtFieldName: []interface{}{
							map[string]interface{}{
								backupPolicyTimeOfDayHoursFieldName:    10,
								backupPolicyTimeOfDayMinutesFieldName:  10,
								backupPolicyTimeOfDayTimeZoneFieldName: "UTC",
							},
						},
					},
				},
			},
		}
		expected[backupPolicyScheduleFieldName] = expectedSchedule
		flattened := flattenBackupPolicyResource(policy)
		assert.Equal(t, expected, flattened)
	})
}
