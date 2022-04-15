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

package pkg

import (
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	backup "github.com/arangodb-managed/apis/backup/v1"
)

func TestFlattenBackupPolicy(t *testing.T) {
	policy := &backup.BackupPolicy{
		Name:              "test-policy",
		Description:       "test-description",
		DeploymentId:      "123456",
		IsPaused:          true,
		Upload:            true,
		RetentionPeriod:   types.DurationProto(200 * time.Hour),
		EmailNotification: "None",
	}

	expected := map[string]interface{}{
		backupPolicyNameFieldName:              "test-policy",
		backupPolicyDescriptionFieldName:       "test-description",
		backupPolicyDeploymentIDFieldName:      "123456",
		backupPolicyIsPausedFieldName:          true,
		backupPolicyUploadFieldName:            true,
		backupPolicyRetentionPeriodFieldName:   200,
		backupPolictEmailNotificationFieldName: "None",
	}

	t.Run("with hourly schedule", func(tt *testing.T) {
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: hourlySchedule,
			HourlySchedule: &backup.BackupPolicy_HourlySchedule{
				ScheduleEveryIntervalHours: 10,
			},
		}
		policy.Schedule = schedule
		expectedSchedule := []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: hourlySchedule,
				backupPolicyScheduleHourlyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheduleHourlyScheduleIntervalFieldName: 10,
					},
				},
			},
		}
		expected[backupPolicyScheduleFieldName] = expectedSchedule
		flattened := flattenBackupPolicyResource(policy)
		assert.Equal(tt, expected, flattened)
	})
	t.Run("with daily schedule", func(tt *testing.T) {
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: dailySchedule,
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
				backupPolicyScheduleTypeFieldName: dailySchedule,
				backupPolicyScheduleDailyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheduleDailyScheduleMondayFieldName:    true,
						backupPolicyScheduleDailyScheduleTuesdayFieldName:   false,
						backupPolicyScheduleDailyScheduleWednesdayFieldName: false,
						backupPolicyScheduleDailyScheduleThursdayFieldName:  true,
						backupPolicyScheduleDailyScheduleFridayFieldName:    false,
						backupPolicyScheduleDailyScheduleSaturdayFieldName:  false,
						backupPolicyScheduleDailyScheduleSundayFieldName:    false,
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
		assert.Equal(tt, expected, flattened)
	})
	t.Run("with monthly schedule", func(tt *testing.T) {
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: monthlySchedule,
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
				backupPolicyScheduleTypeFieldName: monthlySchedule,
				backupPolicyScheduleMonthlyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheduleMonthlyScheduleDayOfMonthScheduleFieldName: 30,
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
		assert.Equal(tt, expected, flattened)
	})
}

func TestExpandBackupPolicy(t *testing.T) {
	raw := map[string]interface{}{
		backupPolicyNameFieldName:              "test-policy",
		backupPolicyDescriptionFieldName:       "test-description",
		backupPolicyDeploymentIDFieldName:      "123456",
		backupPolicyIsPausedFieldName:          true,
		backupPolicyUploadFieldName:            true,
		backupPolicyRetentionPeriodFieldName:   200,
		backupPolictEmailNotificationFieldName: "None",
	}
	expected := &backup.BackupPolicy{
		Name:              "test-policy",
		Description:       "test-description",
		DeploymentId:      "123456",
		IsPaused:          true,
		Upload:            true,
		RetentionPeriod:   types.DurationProto(200 * time.Hour),
		EmailNotification: "None",
	}
	t.Run("test hourly schedule", func(tt *testing.T) {
		rawSchedule := []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: hourlySchedule,
				backupPolicyScheduleHourlyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheduleHourlyScheduleIntervalFieldName: 6,
					},
				},
			},
		}
		raw[backupPolicyScheduleFieldName] = rawSchedule
		s := resourceBackupPolicy().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		policy, err := expandBackupPolicyResource(resourceData)
		assert.NoError(t, err)
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: hourlySchedule,
			HourlySchedule: &backup.BackupPolicy_HourlySchedule{
				ScheduleEveryIntervalHours: 6,
			},
		}
		expected.Schedule = schedule
		assert.Equal(t, expected, policy)
	})
	t.Run("test daily schedule", func(tt *testing.T) {
		rawSchedule := []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: dailySchedule,
				backupPolicyScheduleDailyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheduleDailyScheduleMondayFieldName:    true,
						backupPolicyScheduleDailyScheduleTuesdayFieldName:   true,
						backupPolicyScheduleDailyScheduleWednesdayFieldName: true,
						backupPolicyScheduleDailyScheduleThursdayFieldName:  true,
						backupPolicyScheduleDailyScheduleFridayFieldName:    true,
						backupPolicyScheduleDailyScheduleSaturdayFieldName:  true,
						backupPolicyScheduleDailyScheduleSundayFieldName:    true,
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
		raw[backupPolicyScheduleFieldName] = rawSchedule
		s := resourceBackupPolicy().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		policy, err := expandBackupPolicyResource(resourceData)
		assert.NoError(t, err)
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: dailySchedule,
			DailySchedule: &backup.BackupPolicy_DailySchedule{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
				ScheduleAt: &backup.TimeOfDay{
					Hours:    10,
					Minutes:  10,
					TimeZone: "UTC",
				},
			},
		}
		expected.Schedule = schedule
		assert.Equal(t, expected, policy)
	})
	t.Run("test monthly schedule", func(tt *testing.T) {
		rawSchedule := []interface{}{
			map[string]interface{}{
				backupPolicyScheduleTypeFieldName: monthlySchedule,
				backupPolicyScheduleMonthlyScheduleFieldName: []interface{}{
					map[string]interface{}{
						backupPolicyScheduleMonthlyScheduleDayOfMonthScheduleFieldName: 30,
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
		raw[backupPolicyScheduleFieldName] = rawSchedule
		s := resourceBackupPolicy().Schema
		resourceData := schema.TestResourceDataRaw(t, s, raw)
		policy, err := expandBackupPolicyResource(resourceData)
		assert.NoError(t, err)
		schedule := &backup.BackupPolicy_Schedule{
			ScheduleType: monthlySchedule,
			MonthlySchedule: &backup.BackupPolicy_MonthlySchedule{
				ScheduleAt: &backup.TimeOfDay{
					Hours:    10,
					Minutes:  10,
					TimeZone: "UTC",
				},
				DayOfMonth: 30,
			},
		}
		expected.Schedule = schedule
		assert.Equal(t, expected, policy)
	})
}
