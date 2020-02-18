//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/types"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// Backup field names
	// Main
	backupPolicyNameFieldName         = "name"
	backupPolicyDescriptionFieldName  = "description"
	backupPolicyIsPausedFieldName     = "is_paused"
	backupPolicyScheduleFieldName     = "schedule"
	backupPolicyDeploymentIDFieldName = "deployment_id"
	// Schedule
	backupPolicyScheduleTypeFieldName = "type"
	// Hourly
	backupPolicyScheudleHourlyScheduleFieldName         = "hourly"
	backupPolicyScheudleHourlyScheduleIntervalFieldName = "interval"
	// Daily
	backupPolicyScheudleDailyScheduleFieldName          = "daily"
	backupPolicyScheudleDailyScheduleMondayFieldName    = "monday"
	backupPolicyScheudleDailyScheduleTuesdayFieldName   = "tuesday"
	backupPolicyScheudleDailyScheduleWednesdayFieldName = "wednesday"
	backupPolicyScheudleDailyScheduleThursdayFieldName  = "thursday"
	backupPolicyScheudleDailyScheduleFridayFieldName    = "friday"
	backupPolicyScheudleDailyScheduleSaturdayFieldName  = "saturday"
	backupPolicyScheudleDailyScheduleSundayFieldName    = "sunday"
	// Monthly
	backupPolicyScheudleMonthlyScheduleFieldName                   = "monthly"
	backupPolicyScheudleMonthlyScheduleDayOfMonthScheduleFieldName = "day_of_month"
	// Details
	backupPolicyUploadFieldName            = "upload"
	backupPolicyRetentionPeriodFieldName   = "retention_period"
	backupPolictEmailNotificationFeidlName = "email_notification"
	// TimeOfDay
	backupPolicyTimeOfDayScheduleAtFieldName = "schedule_at"
	backupPolicyTimeOfDayHoursFieldName      = "hours"
	backupPolicyTimeOfDayMinutesFieldName    = "minutes"
	backupPolicyTimeOfDayTimeZoneFieldName   = "timezone"
)

func resourceBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackupPolicyCreate,
		Read:   resourceBackupPolicyRead,
		Update: resourceBackupPolicyUpdate,
		Delete: resourceBackupPolicyDelete,

		Schema: map[string]*schema.Schema{
			backupPolicyNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			backupPolicyDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			backupPolicyIsPausedFieldName: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			backupPolicyUploadFieldName: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			backupPolicyDeploymentIDFieldName: {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// This is a read-only field.
					return true
				},
			},
			backupPolicyRetentionPeriodFieldName: {
				Type: schema.TypeInt,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "0"
				},
				Optional: true,
			},
			backupPolictEmailNotificationFeidlName: {
				Type:     schema.TypeString,
				Required: true,
			},
			backupPolicyScheduleFieldName: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						backupPolicyScheduleTypeFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
						// Hourly
						backupPolicyScheudleHourlyScheduleFieldName: {
							Type: schema.TypeList,
							// Not supported as of now. Enable this check once this issue is fixed:
							// https://github.com/hashicorp/terraform-plugin-sdk/issues/71
							//ConflictsWith: []string{
							//	backupPolicyScheudleDailyScheduleFieldName,
							//	backupPolicyScheudleMonthlyScheduleFieldName,
							//},
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									backupPolicyScheudleHourlyScheduleIntervalFieldName: {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						// Daily
						backupPolicyScheudleDailyScheduleFieldName: {
							Type: schema.TypeList,
							// Not supported as of now. Enable this check once this issue is fixed:
							// https://github.com/hashicorp/terraform-plugin-sdk/issues/71
							//ConflictsWith: []string{
							//	backupPolicyScheudleHourlyScheduleFieldName,
							//	backupPolicyScheudleMonthlyScheduleFieldName,
							//},
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									backupPolicyScheudleDailyScheduleMondayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheudleDailyScheduleTuesdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheudleDailyScheduleWednesdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheudleDailyScheduleThursdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheudleDailyScheduleFridayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheudleDailyScheduleSaturdayFieldName: {
										Type:     schema.TypeBool,
										Optional: true,
									},
									backupPolicyScheudleDailyScheduleSundayFieldName: {
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
						backupPolicyScheudleMonthlyScheduleFieldName: {
							Type: schema.TypeList,
							// Not supported as of now. Enable this check once this issue is fixed:
							// https://github.com/hashicorp/terraform-plugin-sdk/issues/71
							//ConflictsWith: []string{
							//	backupPolicyScheudleDailyScheduleFieldName,
							//	backupPolicyScheudleHourlyScheduleFieldName,
							//},
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									backupPolicyScheudleMonthlyScheduleDayOfMonthScheduleFieldName: {
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
		},
	}
}

func resourceBackupPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBackupPolicyRead(d, m)
}

func resourceBackupPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	policy, err := backupc.GetBackupPolicy(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("backup-policy-id", d.Id()).Msg("Failed to find backup policy")
		d.SetId("")
		return err
	}
	if policy == nil {
		client.log.Error().Err(err).Str("backup-policy-id", d.Id()).Msg("Failed to find backup policy")
		d.SetId("")
		return nil
	}

	for k, v := range flattenBackupPolicyResource(policy) {
		if _, ok := d.GetOk(k); ok {
			if err := d.Set(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func flattenBackupPolicyResource(policy *backup.BackupPolicy) map[string]interface{} {
	schedule := flattenSchedule(policy.GetSchedule())
	ret := map[string]interface{}{
		backupPolicyNameFieldName:              policy.GetName(),
		backupPolicyDescriptionFieldName:       policy.GetDescription(),
		backupPolicyDeploymentIDFieldName:      policy.GetDeploymentId(),
		backupPolicyIsPausedFieldName:          policy.GetIsPaused(),
		backupPolicyUploadFieldName:            policy.GetUpload(),
		backupPolictEmailNotificationFeidlName: policy.GetEmailNotification(),
		backupPolicyScheduleFieldName:          schedule,
	}
	if policy.GetRetentionPeriod() != nil {
		seconds := policy.GetRetentionPeriod().GetSeconds()
		days := seconds / (24 * 60 * 60)
		ret[backupPolicyRetentionPeriodFieldName] = int(days)
	}
	return ret
}

func flattenSchedule(policy *backup.BackupPolicy_Schedule) []interface{} {
	schedule := make(map[string]interface{})
	if policy.GetHourlySchedule() != nil {
		schedule[backupPolicyScheudleHourlyScheduleFieldName] = flattenScheduleHourly(policy.GetHourlySchedule())
	}
	if policy.GetDailySchedule() != nil {
		schedule[backupPolicyScheudleDailyScheduleFieldName] = flattenScheduleDaily(policy.GetDailySchedule())
	}
	if policy.GetMonthlySchedule() != nil {
		schedule[backupPolicyScheudleMonthlyScheduleFieldName] = flattenScheduleMonthly(policy.GetMonthlySchedule())
	}
	schedule[backupPolicyScheduleTypeFieldName] = policy.GetScheduleType()
	return []interface{}{
		schedule,
	}
}

func flattenScheduleHourly(policy *backup.BackupPolicy_HourlySchedule) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyScheudleHourlyScheduleIntervalFieldName: int(policy.GetScheduleEveryIntervalHours()),
		},
	}
}

func flattenScheduleDaily(policy *backup.BackupPolicy_DailySchedule) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyScheudleDailyScheduleMondayFieldName:    policy.GetMonday(),
			backupPolicyScheudleDailyScheduleTuesdayFieldName:   policy.GetTuesday(),
			backupPolicyScheudleDailyScheduleWednesdayFieldName: policy.GetWednesday(),
			backupPolicyScheudleDailyScheduleThursdayFieldName:  policy.GetThursday(),
			backupPolicyScheudleDailyScheduleFridayFieldName:    policy.GetFriday(),
			backupPolicyScheudleDailyScheduleSaturdayFieldName:  policy.GetSaturday(),
			backupPolicyScheudleDailyScheduleSundayFieldName:    policy.GetSunday(),
			backupPolicyTimeOfDayScheduleAtFieldName:            flattenTimeOfDay(policy.GetScheduleAt()),
		},
	}
}

func flattenScheduleMonthly(policy *backup.BackupPolicy_MonthlySchedule) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyScheudleMonthlyScheduleDayOfMonthScheduleFieldName: int(policy.GetDayOfMonth()),
			backupPolicyTimeOfDayScheduleAtFieldName:                       flattenTimeOfDay(policy.GetScheduleAt()),
		},
	}
}

func flattenTimeOfDay(day *backup.TimeOfDay) []interface{} {
	return []interface{}{
		map[string]interface{}{
			backupPolicyTimeOfDayHoursFieldName:    int(day.GetHours()),
			backupPolicyTimeOfDayMinutesFieldName:  int(day.GetMinutes()),
			backupPolicyTimeOfDayTimeZoneFieldName: day.GetTimeZone(),
		},
	}
}

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
		ret.RetentionPeriod = types.DurationProto(time.Duration(v.(int)))
	}
	if v, ok := d.GetOk(backupPolictEmailNotificationFeidlName); ok {
		ret.EmailNotification = v.(string)
	}
	if v, ok := d.GetOk(backupPolicyScheduleFieldName); ok {
		expandBackupPolicySchedule(v.([]interface{}))
	}
	return ret, nil
}

func expandBackupPolicySchedule(s []interface{}) *backup.BackupPolicy_Schedule {
	ret := &backup.BackupPolicy_Schedule{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheduleTypeFieldName]; ok {
			ret.ScheduleType = i.(string)
		}
		if i, ok := item[backupPolicyScheudleHourlyScheduleFieldName]; ok {
			ret.HourlySchedule = expandHourlySchedule(i.([]interface{}))
		}
		if i, ok := item[backupPolicyScheudleDailyScheduleFieldName]; ok {
			ret.DailySchedule = expandDailySchedule(i.([]interface{}))
		}
		if i, ok := item[backupPolicyScheudleMonthlyScheduleFieldName]; ok {
			ret.MonthlySchedule = expandMonthlySchedule(i.([]interface{}))
		}
	}
	return ret
}

func expandMonthlySchedule(s []interface{}) *backup.BackupPolicy_MonthlySchedule {
	ret := &backup.BackupPolicy_MonthlySchedule{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheudleMonthlyScheduleDayOfMonthScheduleFieldName]; ok {
			ret.DayOfMonth = int32(i.(int))
		}
		if i, ok := item[backupPolicyTimeOfDayScheduleAtFieldName]; ok {
			ret.ScheduleAt = expandTimeOfDay(i.([]interface{}))
		}
	}
	return ret
}

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

func expandDailySchedule(s []interface{}) *backup.BackupPolicy_DailySchedule {
	ret := &backup.BackupPolicy_DailySchedule{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheudleDailyScheduleMondayFieldName]; ok {
			ret.Monday = i.(bool)
		}
		if i, ok := item[backupPolicyScheudleDailyScheduleTuesdayFieldName]; ok {
			ret.Tuesday = i.(bool)
		}
		if i, ok := item[backupPolicyScheudleDailyScheduleWednesdayFieldName]; ok {
			ret.Wednesday = i.(bool)
		}
		if i, ok := item[backupPolicyScheudleDailyScheduleThursdayFieldName]; ok {
			ret.Thursday = i.(bool)
		}
		if i, ok := item[backupPolicyScheudleDailyScheduleFridayFieldName]; ok {
			ret.Friday = i.(bool)
		}
		if i, ok := item[backupPolicyScheudleDailyScheduleSaturdayFieldName]; ok {
			ret.Saturday = i.(bool)
		}
		if i, ok := item[backupPolicyScheudleDailyScheduleSundayFieldName]; ok {
			ret.Sunday = i.(bool)
		}
		if i, ok := item[backupPolicyTimeOfDayScheduleAtFieldName]; ok {
			ret.ScheduleAt = expandTimeOfDay(i.([]interface{}))
		}
	}
	return ret
}

func expandHourlySchedule(s []interface{}) *backup.BackupPolicy_HourlySchedule {
	ret := &backup.BackupPolicy_HourlySchedule{}
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[backupPolicyScheudleHourlyScheduleIntervalFieldName]; ok {
			ret.ScheduleEveryIntervalHours = int32(i.(int))
		}
	}
	return ret
}

func resourceBackupPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	expandedPolicy, err := expandBackupPolicyResource(d)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to expand on policy")
		return err
	}
	if b, err := backupc.CreateBackupPolicy(client.ctxWithToken, expandedPolicy); err != nil {
		client.log.Error().Err(err).Msg("Failed to create backup policy")
		return err
	} else {
		d.SetId(b.GetId())
	}
	return resourceBackupPolicyRead(d, m)
}

func resourceBackupPolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	if _, err := backupc.DeleteBackupPolicy(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("backup-policy-id", d.Id()).Msg("Failed to delete backup policy")
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}
