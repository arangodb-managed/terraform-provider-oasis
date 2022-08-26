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
	"errors"
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
	backupNameFieldName         = "name"
	backupDescriptionFieldName  = "description"
	backupURLFieldName          = "url"
	backupPolicyIDFieldName     = "backup_policy_id"
	backupDeploymentIDFieldName = "deployment_id"
	backupUploadFieldName       = "upload"
	backupAutoDeleteAtFieldName = "auto_deleted_at"
)

// resourceBackup defines a Backup Oasis resource.
func resourceBackup() *schema.Resource {
	return &schema.Resource{
		Description:   "Oasis Backup Resource",
		CreateContext: resourceBackupCreate,
		ReadContext:   resourceBackupRead,
		UpdateContext: resourceBackupUpdate,
		DeleteContext: resourceBackupDelete,
		Schema: map[string]*schema.Schema{
			backupNameFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Backup Resource Backup Name field",
				Required:    true,
			},
			backupDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Backup Resource Backup Description field",
				Optional:    true,
			},
			backupUploadFieldName: {
				Type:        schema.TypeBool,
				Description: "Oasis Backup Resource Backup Upload field",
				Optional:    true,
				Default:     false,
			},
			backupDeploymentIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Backup Resource Backup Deployment ID field",
				Required:    true,
			},
			backupURLFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Backup Resource Backup URL field",
				Computed:    true,
			},
			backupPolicyIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Backup Resource Backup Policy ID field",
				Optional:    true,
			},
			backupAutoDeleteAtFieldName: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Oasis Backup Resource Backup Auto Delete At field",
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					days := v.(int)
					var errs []error
					if days < 0 || days > 31 {
						errs = append(errs, errors.New("auto_delete_at: must be within range 1-31"))
					}
					return nil, errs
				},
			},
		},
	}
}

// resourceBackupRead will gather information from the Terraform store and display it accordingly.
func resourceBackupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	backup, err := backupc.GetBackup(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil || backup == nil {
		client.log.Error().Err(err).Str("backup-id", d.Id()).Msg("Failed to find backup")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenBackupResource(backup) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// expandBackupResource will take a Terraform flat map schema data and turn it into an Oasis Backup.
func expandBackupResource(d *schema.ResourceData) (*backup.Backup, error) {
	ret := &backup.Backup{}
	if v, ok := d.GetOk(backupNameFieldName); ok {
		ret.Name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", backupNameFieldName)
	}
	if v, ok := d.GetOk(backupDescriptionFieldName); ok {
		ret.Description = v.(string)
	}
	if v, ok := d.GetOk(backupUploadFieldName); ok {
		ret.Upload = v.(bool)
	}
	if v, ok := d.GetOk(backupDeploymentIDFieldName); ok {
		ret.DeploymentId = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", backupDeploymentIDFieldName)
	}
	if v, ok := d.GetOk(backupAutoDeleteAtFieldName); ok {
		autoDeleteAt, err := types.TimestampProto(time.Now().AddDate(0, 0, v.(int))) // add n days for backup auto deletion
		if err != nil {
			return nil, fmt.Errorf("unable to parse time for auto delete backup")
		}
		ret.AutoDeletedAt = autoDeleteAt
	}

	return ret, nil
}

// resourceBackupCreate will take the schema data from the Terraform config file and call the Oasis client
// to initiate a create procedure for a Backup. It will call helper methods to construct the necessary data
// in order to create this object.
func resourceBackupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	expandedBackup, err := expandBackupResource(d)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to expand on backup")
		return diag.FromErr(err)
	}

	// Pre-check for the given deployment
	datac := data.NewDataServiceClient(client.conn)
	if _, err := datac.GetDeployment(client.ctxWithToken, &common.IDOptions{Id: expandedBackup.DeploymentId}); err != nil {
		client.log.Error().Err(err).Str("deployment-id", expandedBackup.DeploymentId).Msg("Deployment with ID not found.")
		return diag.FromErr(err)
	}

	if b, err := backupc.CreateBackup(client.ctxWithToken, expandedBackup); err != nil {
		client.log.Error().Err(err).Msg("Failed to create backup")
		return diag.FromErr(err)
	} else {
		d.SetId(b.GetId())
	}

	return resourceBackupRead(ctx, d, m)
}

// resourceBackupDelete will delete a given resource based on the calculated ID.
func resourceBackupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	if _, err := backupc.DeleteBackup(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("backup-id", d.Id()).Msg("Failed to delete backup")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceBackupUpdate will take a resource diff and apply changes accordingly if there are any.
func resourceBackupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	backupc := backup.NewBackupServiceClient(client.conn)
	backup, err := backupc.GetBackup(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find backup")
		d.SetId("")
		return diag.FromErr(err)
	}
	// Main fields
	if d.HasChange(backupNameFieldName) {
		backup.Name = d.Get(backupNameFieldName).(string)
	}
	if d.HasChange(backupDescriptionFieldName) {
		backup.Description = d.Get(backupDescriptionFieldName).(string)
	}
	if d.HasChange(backupUploadFieldName) {
		backup.Upload = d.Get(backupUploadFieldName).(bool)
	}
	if d.HasChange(backupAutoDeleteAtFieldName) {
		updatedAutoDeleteAt, err := types.TimestampProto(time.Now().AddDate(0, 0, d.Get(backupAutoDeleteAtFieldName).(int)))
		if err != nil {
			return diag.FromErr(errors.New("unable to parse time for auto delete backup"))
		}
		backup.AutoDeletedAt = updatedAutoDeleteAt
	}

	res, err := backupc.UpdateBackup(client.ctxWithToken, backup)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update backup")
		return diag.FromErr(err)
	} else {
		d.SetId(res.GetId())
	}
	return resourceBackupRead(ctx, d, m)
}

// flattenBackupResource will take a Backup object and turn it into a flat map for terraform digestion.
func flattenBackupResource(backup *backup.Backup) map[string]interface{} {
	return map[string]interface{}{
		backupNameFieldName:         backup.GetName(),
		backupDescriptionFieldName:  backup.GetDescription(),
		backupURLFieldName:          backup.GetUrl(),
		backupPolicyIDFieldName:     backup.GetBackupPolicyId(),
		backupDeploymentIDFieldName: backup.GetDeploymentId(),
		backupRegionID:              backup.GetRegionId(),
	}
}
