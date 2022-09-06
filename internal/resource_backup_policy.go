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
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
)

const (
	// Backup field names
	// Main
	backupPolicyNameFieldName         = "name"
	backupPolicyDescriptionFieldName  = "description"
	backupPolicyIsPausedFieldName     = "is_paused"
	backupPolicyScheduleFieldName     = "schedule"
	backupPolicyDeploymentIDFieldName = "deployment_id"
	backupPolicyLockedFieldName       = "locked"
	// Schedule
	backupPolicyScheduleTypeFieldName = "type"
	// Hourly
	backupPolicyScheduleHourlyScheduleFieldName         = "hourly"
	backupPolicyScheduleHourlyScheduleIntervalFieldName = "interval"
	// Daily
	backupPolicyScheduleDailyScheduleFieldName          = "daily"
	backupPolicyScheduleDailyScheduleMondayFieldName    = "monday"
	backupPolicyScheduleDailyScheduleTuesdayFieldName   = "tuesday"
	backupPolicyScheduleDailyScheduleWednesdayFieldName = "wednesday"
	backupPolicyScheduleDailyScheduleThursdayFieldName  = "thursday"
	backupPolicyScheduleDailyScheduleFridayFieldName    = "friday"
	backupPolicyScheduleDailyScheduleSaturdayFieldName  = "saturday"
	backupPolicyScheduleDailyScheduleSundayFieldName    = "sunday"
	// Monthly
	backupPolicyScheduleMonthlyScheduleFieldName                   = "monthly"
	backupPolicyScheduleMonthlyScheduleDayOfMonthScheduleFieldName = "day_of_month"
	// Details
	backupPolicyUploadFieldName            = "upload"
	backupPolicyRetentionPeriodFieldName   = "retention_period_hour"
	backupPolictEmailNotificationFieldName = "email_notification"
	// TimeOfDay
	backupPolicyTimeOfDayScheduleAtFieldName = "schedule_at"
	backupPolicyTimeOfDayHoursFieldName      = "hours"
	backupPolicyTimeOfDayMinutesFieldName    = "minutes"
	backupPolicyTimeOfDayTimeZoneFieldName   = "timezone"

	// Schedule Types
	hourlySchedule  = "Hourly"
	dailySchedule   = "Daily"
	monthlySchedule = "Monthly"

	// Additional region identifiers where backup should be cloned
	backupPolicyAdditionalRegionIDs = "additional_region_ids"
)

// resourceBackupPolicy defines a BackupPolicy oasis resource.
func resourceBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Backup Policy Resource",

		CreateContext: resourceBackupPolicyCreate,
		ReadContext:   resourceBackupPolicyRead,
		UpdateContext: resourceBackupPolicyUpdate,
		DeleteContext: resourceBackupPolicyDelete,

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			o, n := diff.GetChange(backupPolicyDeploymentIDFieldName)
			if o != "" && o != n {
				return fmt.Errorf("Cannot change deployment ID once it has been set.")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			backupPolicyNameFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Policy Resource Backup Policy Name field",
				Required:    true,
			},
			backupPolicyDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Policy Resource Backup Policy Description field",
				Optional:    true,
			},
			backupPolicyIsPausedFieldName: {
				Type:        schema.TypeBool,
				Description: "Backup Policy Resource Backup Policy Is Paused field",
				Optional:    true,
			},
			backupPolicyUploadFieldName: {
				Type:        schema.TypeBool,
				Description: "Backup Policy Resource Backup Policy Upload field",
				Optional:    true,
			},
			backupPolicyDeploymentIDFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Policy Resource Backup Policy Deployment ID field",
				Required:    true,
			},
			backupPolicyRetentionPeriodFieldName: {
				Type:        schema.TypeInt,
				Description: "Backup Policy Resource Backup Policy Retention Period field",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "0"
				},
				Optional: true,
			},
			backupPolictEmailNotificationFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Policy Resource Backup Policy Email Notification field",
				Required:    true,
			},
			backupPolicyAdditionalRegionIDs: {
				Type:        schema.TypeList,
				Description: "Backup Policy Additional Region Identifiers where backup should be cloned",
				Optional:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			backupPolicyScheduleFieldName: {
				Type:        schema.TypeList,
				Description: "Backup Policy Resource Backup Policy Schedule field",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						backupPolicyScheduleTypeFieldName: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Schedule type should be one of the following string: \"Hourly|Daily|Monthly\"",
						},
						// Hourly
						backupPolicyScheduleHourlyScheduleFieldName: {
							Type: schema.TypeList,
							// Not supported as of now. Enable this check once this issue is fixed:
							// https://github.com/hashicorp/terraform-plugin-sdk/issues/71
							//ConflictsWith: []string{
							//	backupPolicyScheduleDailyScheduleFieldName,
							//	backupPolicyScheduleMonthlyScheduleFieldName,
							//},
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									backupPolicyScheduleHourlyScheduleIntervalFieldName: {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						// Daily
						backupPolicyScheduleDailyScheduleFieldName: {
							Type: schema.TypeList,
							// Not supported as of now. Enable this check once this issue is fixed:
							// https://github.com/hashicorp/terraform-plugin-sdk/issues/71
							//ConflictsWith: []string{
							//	backupPolicyScheduleHourlyScheduleFieldName,
							//	backupPolicyScheduleMonthlyScheduleFieldName,
							//},
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									backupPolicyScheduleDailyScheduleMondayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheduleDailyScheduleTuesdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheduleDailyScheduleWednesdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheduleDailyScheduleThursdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheduleDailyScheduleFridayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheduleDailyScheduleSaturdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheduleDailyScheduleSundayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyTimeOfDayScheduleAtFieldName: {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												backupPolicyTimeOfDayHoursFieldName: {
													Type:     schema.TypeInt,
													Optional: true,
												},
												backupPolicyTimeOfDayMinutesFieldName: {
													Type:     schema.TypeInt,
													Optional: true,
												},
												backupPolicyTimeOfDayTimeZoneFieldName: {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						// Monthly
						backupPolicyScheduleMonthlyScheduleFieldName: {
							Type: schema.TypeList,
							// Not supported as of now. Enable this check once this issue is fixed:
							// https://github.com/hashicorp/terraform-plugin-sdk/issues/71
							//ConflictsWith: []string{
							//	backupPolicyScheduleDailyScheduleFieldName,
							//	backupPolicyScheduleHourlyScheduleFieldName,
							//},
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									backupPolicyScheduleMonthlyScheduleDayOfMonthScheduleFieldName: {
										Type:     schema.TypeInt,
										Optional: true,
									},
									backupPolicyTimeOfDayScheduleAtFieldName: {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												backupPolicyTimeOfDayHoursFieldName: {
													Type:     schema.TypeInt,
													Optional: true,
												},
												backupPolicyTimeOfDayMinutesFieldName: {
													Type:     schema.TypeInt,
													Optional: true,
												},
												backupPolicyTimeOfDayTimeZoneFieldName: {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			backupPolicyLockedFieldName: {
				Type:        schema.TypeBool,
				Description: "Backup Policy Resource Backup Policy Locked field",
				Optional:    true,
			},
		},
	}
}

// resourceBackupPolicyUpdate will take a resource diff and apply changes accordingly if there are any.
func resourceBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	backupc := backup.NewBackupServiceClient(client.conn)
	policy, err := backupc.GetBackupPolicy(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find backup policy")
		d.SetId("")
		return diag.FromErr(err)
	}
	// Main fields
	if d.HasChange(backupPolicyNameFieldName) {
		policy.Name = d.Get(backupPolicyNameFieldName).(string)
	}
	if d.HasChange(backupPolicyDescriptionFieldName) {
		policy.Description = d.Get(backupPolicyDescriptionFieldName).(string)
	}
	if d.HasChange(backupPolicyIsPausedFieldName) {
		policy.IsPaused = d.Get(backupPolicyIsPausedFieldName).(bool)
	}
	if d.HasChange(backupPolicyUploadFieldName) {
		policy.Upload = d.Get(backupPolicyUploadFieldName).(bool)
	}
	if d.HasChange(backupPolicyRetentionPeriodFieldName) {
		v := d.Get(backupPolicyRetentionPeriodFieldName)
		policy.RetentionPeriod = getRetentionPeriod(v)
	}
	if d.HasChange(backupPolictEmailNotificationFieldName) {
		policy.EmailNotification = d.Get(backupPolictEmailNotificationFieldName).(string)
	}
	if d.HasChange(backupPolicyScheduleFieldName) {
		policy.Schedule = expandBackupPolicySchedule(d.Get(backupPolicyScheduleFieldName).([]interface{}))
	}

	// if v, ok := d.GetOk(backupPolicyAdditionalRegionIDs); ok {
	// 	additionalRegionIDs, err := expandAdditionalRegionList(v.([]interface{}))
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// 	policy.AdditionalRegionIds = additionalRegionIDs
	// }
	// Make sure we are sending the right schedule. The coming back schedule can contain invalid
	// field items for different schedule. We make sure here that the right one is sent after an
	// update. This check can be removed once terraform allows conflict checks for list items.
	switch policy.GetSchedule().GetScheduleType() {
	case hourlySchedule:
		policy.Schedule.DailySchedule = nil
		policy.Schedule.MonthlySchedule = nil
	case dailySchedule:
		policy.Schedule.HourlySchedule = nil
		policy.Schedule.MonthlySchedule = nil
	case monthlySchedule:
		policy.Schedule.HourlySchedule = nil
		policy.Schedule.DailySchedule = nil
	}
	if d.HasChange(backupPolicyLockedFieldName) {
		policy.Locked = d.Get(backupPolicyLockedFieldName).(bool)
	}

	res, err := backupc.UpdateBackupPolicy(client.ctxWithToken, policy)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update backup policy")
		return diag.FromErr(err)
	} else {
		d.SetId(res.GetId())
	}
	return resourceBackupPolicyRead(ctx, d, m)
}

// getRetentionPeriod calculates the retention period.
func getRetentionPeriod(v interface{}) *types.Duration {
	// retention period is given in hours
	return types.DurationProto((time.Duration(v.(int)) * 60 * 60) * time.Second)
}

// resourceBackupPolicyRead will gather information from the terraform store and display it accordingly.
func resourceBackupPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	policy, err := backupc.GetBackupPolicy(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("backup-policy-id", d.Id()).Msg("Failed to find backup policy")
		d.SetId("")
		return diag.FromErr(err)
	}
	if policy == nil {
		client.log.Error().Err(err).Str("backup-policy-id", d.Id()).Msg("Failed to find backup policy")
		d.SetId("")
		return nil
	}

	for k, v := range flattenBackupPolicyResource(policy) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// flattenBackupPolicyResource will take a BackupPolicy object and turn it into a flat map for terraform digestion.
func flattenBackupPolicyResource(policy *backup.BackupPolicy) map[string]interface{} {
	schedule := flattenSchedule(policy.GetSchedule())
	ret := map[string]interface{}{
		backupPolicyNameFieldName:              policy.GetName(),
		backupPolicyDescriptionFieldName:       policy.GetDescription(),
		backupPolicyDeploymentIDFieldName:      policy.GetDeploymentId(),
		backupPolicyIsPausedFieldName:          policy.GetIsPaused(),
		backupPolicyUploadFieldName:            policy.GetUpload(),
		backupPolictEmailNotificationFieldName: policy.GetEmailNotification(),
		backupPolicyScheduleFieldName:          schedule,
		backupPolicyLockedFieldName:            policy.GetLocked(),
		backupPolicyAdditionalRegionIDs:        policy.GetAdditionalRegionIds(),
	}
	if policy.GetRetentionPeriod() != nil {
		seconds := policy.GetRetentionPeriod().GetSeconds()
		hours := seconds / (60 * 60)
		ret[backupPolicyRetentionPeriodFieldName] = int(hours)
	}
	return ret
}

// flattenBackupPolicyResource will take a Schedule portion of a BackupPolicy object and turn it into a flat map for terraform digestion.
func flattenSchedule(policy *backup.BackupPolicy_Schedule) []interface{} {
	schedule := make(map[string]interface{})
	schedule[backupPolicyScheduleTypeFieldName] = policy.GetScheduleType()
	switch policy.GetScheduleType() {
	case hourlySchedule:
		schedule[backupPolicyScheduleHourlyScheduleFieldName] = flattenScheduleHourly(policy.GetHourlySchedule())
	case dailySchedule:
		schedule[backupPolicyScheduleDailyScheduleFieldName] = flattenScheduleDaily(policy.GetDailySchedule())
	case monthlySchedule:
		schedule[backupPolicyScheduleMonthlyScheduleFieldName] = flattenScheduleMonthly(policy.GetMonthlySchedule())
	}
	return []interface{}{
		schedule,
	}
}

// flattenScheduleHourly will take the Hourly portion of a schedule an turn it into a flat map.
func flattenScheduleHourly(policy *backup.BackupPolicy_HourlySchedule) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyScheduleHourlyScheduleIntervalFieldName: int(policy.GetScheduleEveryIntervalHours()),
		},
	}
}

// flattenScheduleDaily will take the Daily portion of a schedule an turn it into a flat map.
func flattenScheduleDaily(policy *backup.BackupPolicy_DailySchedule) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyScheduleDailyScheduleMondayFieldName:    policy.GetMonday(),
			backupPolicyScheduleDailyScheduleTuesdayFieldName:   policy.GetTuesday(),
			backupPolicyScheduleDailyScheduleWednesdayFieldName: policy.GetWednesday(),
			backupPolicyScheduleDailyScheduleThursdayFieldName:  policy.GetThursday(),
			backupPolicyScheduleDailyScheduleFridayFieldName:    policy.GetFriday(),
			backupPolicyScheduleDailyScheduleSaturdayFieldName:  policy.GetSaturday(),
			backupPolicyScheduleDailyScheduleSundayFieldName:    policy.GetSunday(),
			backupPolicyTimeOfDayScheduleAtFieldName:            flattenTimeOfDay(policy.GetScheduleAt()),
		},
	}
}

// flattenScheduleMonthly will take the Monthly portion of a schedule an turn it into a flat map.
func flattenScheduleMonthly(policy *backup.BackupPolicy_MonthlySchedule) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyScheduleMonthlyScheduleDayOfMonthScheduleFieldName: int(policy.GetDayOfMonth()),
			backupPolicyTimeOfDayScheduleAtFieldName:                       flattenTimeOfDay(policy.GetScheduleAt()),
		},
	}
}

// flattenTimeOfDay will take the TimeOfDay portion of a schedule an turn it into a flat map.
func flattenTimeOfDay(day *backup.TimeOfDay) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyTimeOfDayHoursFieldName:    int(day.GetHours()),
			backupPolicyTimeOfDayMinutesFieldName:  int(day.GetMinutes()),
			backupPolicyTimeOfDayTimeZoneFieldName: day.GetTimeZone(),
		},
	}
}

// expandBackupPolicyResource will take a terraform flat map schema data and turn it into an Oasis BackupPolicy.
func expandBackupPolicyResource(d *schema.ResourceData) (*backup.BackupPolicy, error) {
	ret := &backup.BackupPolicy{}
	if v, ok := d.GetOk(backupPolicyNameFieldName); ok {
		ret.Name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", backupPolicyNameFieldName)
	}
	if v, ok := d.GetOk(backupPolicyDescriptionFieldName); ok {
		ret.Description = v.(string)
	}
	if v, ok := d.GetOk(backupPolicyIsPausedFieldName); ok {
		ret.IsPaused = v.(bool)
	}
	if v, ok := d.GetOk(backupPolicyUploadFieldName); ok {
		ret.Upload = v.(bool)
	}
	if v, ok := d.GetOk(backupPolicyDeploymentIDFieldName); ok {
		ret.DeploymentId = v.(string)
	}
	if v, ok := d.GetOk(backupPolicyRetentionPeriodFieldName); ok {
		ret.RetentionPeriod = getRetentionPeriod(v)
	}
	if v, ok := d.GetOk(backupPolictEmailNotificationFieldName); ok {
		ret.EmailNotification = v.(string)
	}
	if v, ok := d.GetOk(backupPolicyScheduleFieldName); ok {
		ret.Schedule = expandBackupPolicySchedule(v.([]interface{}))
	}
	if v, ok := d.GetOk(backupPolicyLockedFieldName); ok {
		ret.Locked = v.(bool)
	}
	if v, ok := d.GetOk(backupPolicyAdditionalRegionIDs); ok {
		additionalRegionIDs, err := expandAdditionalRegionList(v.([]interface{}))
		if err != nil {
			return ret, err
		}
		ret.AdditionalRegionIds = additionalRegionIDs
	}
	return ret, nil
}

// expandBackupPolicySchedule will take a terraform flat map schema data and turn it into an Oasis BackupPolicy Schedule.
func expandBackupPolicySchedule(s []interface{}) *backup.BackupPolicy_Schedule {
	ret := &backup.BackupPolicy_Schedule{}
	// First, find the schedule type.
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheduleTypeFieldName]; ok {
			ret.ScheduleType = i.(string)
		}
	}
	// Any other schedules need to be cleared out. This is necessary until
	// terraform allows conflict checks for list items.
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheduleHourlyScheduleFieldName]; ok && ret.ScheduleType == hourlySchedule {
			hourlySchedule := i.([]interface{})
			if len(hourlySchedule) > 0 {
				ret.HourlySchedule = expandHourlySchedule(hourlySchedule)
			}
			ret.DailySchedule = nil
			ret.MonthlySchedule = nil
		}
		if i, ok := item[backupPolicyScheduleDailyScheduleFieldName]; ok && ret.ScheduleType == dailySchedule {
			dailySchedule := i.([]interface{})
			if len(dailySchedule) > 0 {
				ret.DailySchedule = expandDailySchedule(dailySchedule)
			}
			ret.HourlySchedule = nil
			ret.MonthlySchedule = nil
		}
		if i, ok := item[backupPolicyScheduleMonthlyScheduleFieldName]; ok && ret.ScheduleType == monthlySchedule {
			monthlySchedule := i.([]interface{})
			if len(monthlySchedule) > 0 {
				ret.MonthlySchedule = expandMonthlySchedule(monthlySchedule)
			}
			ret.DailySchedule = nil
			ret.HourlySchedule = nil
		}
	}
	return ret
}

// expandMonthlySchedule will take a terraform flat map schema data and decipher the monthly schedule from it.
func expandMonthlySchedule(s []interface{}) *backup.BackupPolicy_MonthlySchedule {
	ret := &backup.BackupPolicy_MonthlySchedule{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheduleMonthlyScheduleDayOfMonthScheduleFieldName]; ok {
			ret.DayOfMonth = int32(i.(int))
		}
		if i, ok := item[backupPolicyTimeOfDayScheduleAtFieldName]; ok {
			ret.ScheduleAt = expandTimeOfDay(i.([]interface{}))
		}
	}
	return ret
}

// expandTimeOfDay will take a terraform flat map schema data and decipher the time of day schedule from it.
func expandTimeOfDay(s []interface{}) *backup.TimeOfDay {
	ret := &backup.TimeOfDay{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyTimeOfDayHoursFieldName]; ok {
			ret.Hours = int32(i.(int))
		}
		if i, ok := item[backupPolicyTimeOfDayMinutesFieldName]; ok {
			ret.Minutes = int32(i.(int))
		}
		if i, ok := item[backupPolicyTimeOfDayTimeZoneFieldName]; ok {
			ret.TimeZone = i.(string)
		}
	}
	return ret
}

// expandDailySchedule will take a terraform flat map schema data and decipher the daily schedule from it.
func expandDailySchedule(s []interface{}) *backup.BackupPolicy_DailySchedule {
	ret := &backup.BackupPolicy_DailySchedule{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheduleDailyScheduleMondayFieldName]; ok {
			ret.Monday = i.(bool)
		}
		if i, ok := item[backupPolicyScheduleDailyScheduleTuesdayFieldName]; ok {
			ret.Tuesday = i.(bool)
		}
		if i, ok := item[backupPolicyScheduleDailyScheduleWednesdayFieldName]; ok {
			ret.Wednesday = i.(bool)
		}
		if i, ok := item[backupPolicyScheduleDailyScheduleThursdayFieldName]; ok {
			ret.Thursday = i.(bool)
		}
		if i, ok := item[backupPolicyScheduleDailyScheduleFridayFieldName]; ok {
			ret.Friday = i.(bool)
		}
		if i, ok := item[backupPolicyScheduleDailyScheduleSaturdayFieldName]; ok {
			ret.Saturday = i.(bool)
		}
		if i, ok := item[backupPolicyScheduleDailyScheduleSundayFieldName]; ok {
			ret.Sunday = i.(bool)
		}
		if i, ok := item[backupPolicyTimeOfDayScheduleAtFieldName]; ok {
			ret.ScheduleAt = expandTimeOfDay(i.([]interface{}))
		}
	}
	return ret
}

// expandHourlySchedule will take a terraform flat map schema data and decipher the hourly schedule from it.
func expandHourlySchedule(s []interface{}) *backup.BackupPolicy_HourlySchedule {
	ret := &backup.BackupPolicy_HourlySchedule{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheduleHourlyScheduleIntervalFieldName]; ok {
			ret.ScheduleEveryIntervalHours = int32(i.(int))
		}
	}
	return ret
}

// resourceBackupPolicyCreate will take the schema data from the terraform config file and call the oasis client
// to initiate a create procedure for a BackupPolicy. It will call helper methods to construct the necessary data
// in order to create this object.
func resourceBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	expandedPolicy, err := expandBackupPolicyResource(d)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to expand on policy")
		return diag.FromErr(err)
	}
	// Pre-check for the given deployment
	datac := data.NewDataServiceClient(client.conn)
	if _, err := datac.GetDeployment(client.ctxWithToken, &common.IDOptions{Id: expandedPolicy.DeploymentId}); err != nil {
		client.log.Error().Err(err).Str("deployment-id", expandedPolicy.DeploymentId).Msg("Deployment with ID not found.")
		return diag.FromErr(err)
	}
	if b, err := backupc.CreateBackupPolicy(client.ctxWithToken, expandedPolicy); err != nil {
		client.log.Error().Err(err).Msg("Failed to create backup policy")
		return diag.FromErr(err)
	} else {
		d.SetId(b.GetId())
	}
	return resourceBackupPolicyRead(ctx, d, m)
}

// resourceBackupPolicyDelete will delete a given resource based on the calculated ID.
func resourceBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	if _, err := backupc.DeleteBackupPolicy(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("backup-policy-id", d.Id()).Msg("Failed to delete backup policy")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// expandAdditionalRegionList creates a string list of items from an interface slice. It also
// verifies if a given string item is empty or not. In case it's empty, an error is thrown.
func expandAdditionalRegionList(list []interface{}) ([]string, error) {
	additionalRegionIDs := []string{}
	for _, v := range list {
		if v, ok := v.(string); ok {
			if v == "" {
				continue
			}
			additionalRegionIDs = append(additionalRegionIDs, v)
		}
	}
	return additionalRegionIDs, nil
}
