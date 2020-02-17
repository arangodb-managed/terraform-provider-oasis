//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// Backup field names
	// Main
	backupPolicyNameFieldName        = "name"
	backupPolicyDescriptionFieldName = "description"
	backupPolicyIsPausedFieldName    = "is_paused"
	backupPolicyScheduleFieldName    = "schedule"
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
		backupPolicyIsPausedFieldName:          policy.GetIsPaused(),
		backupPolicyUploadFieldName:            policy.GetUpload(),
		backupPolictEmailNotificationFeidlName: policy.GetEmailNotification(),
		backupPolicyScheduleFieldName:          schedule,
	}
	if policy.GetRetentionPeriod() != nil {
		seconds := policy.GetRetentionPeriod().GetSeconds()
		days := seconds / (24 * 60 * 60)
		ret[backupPolicyRetentionPeriodFieldName] = days
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

func resourceBackupPolicyCreate(d *schema.ResourceData, m interface{}) error {
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
