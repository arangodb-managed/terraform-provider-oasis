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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
)

const (
	// Backup data source field names
	backupDataSourceIdFieldName           = "id"
	backupDataSourceNameFieldName         = "name"
	backupDataSourceDescriptionFieldName  = "description"
	backupDataSourceURLFieldName          = "url"
	backupDataSourcePolicyIDFieldName     = "backup_policy_id"
	backupDataSourceDeploymentIDFieldName = "deployment_id"
	backupDataSourceCreatedAtFieldName    = "created_at"
)

// dataSourceOasisBackup defines a Backup datasource terraform type.
func dataSourceOasisBackup() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Backup Data Source",

		ReadContext: dataSourceOasisBackupRead,

		Schema: map[string]*schema.Schema{
			backupDataSourceIdFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Data Source ID field",
				Required:    true,
			},
			backupDataSourceNameFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Data Source Name field",
				Optional:    true,
			},
			backupDataSourceDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Data Source Description field",
				Optional:    true,
			},
			backupDataSourceURLFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Data Source URL field",
				Optional:    true,
			},
			backupDataSourcePolicyIDFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Data Source Policy ID field",
				Optional:    true,
			},
			backupDataSourceDeploymentIDFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Sata Source Deployment ID field",
				Optional:    true,
			},
			backupDataSourceCreatedAtFieldName: {
				Type:        schema.TypeString,
				Description: "Backup Sata Source Created At field",
				Computed:    true,
			},
		},
	}
}

// dataSourceOasisBackupRead reloads the resource object from the Terraform store.
func dataSourceOasisBackupRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)
	bid := data.Get(backupDataSourceIdFieldName).(string)
	backup, err := backupc.GetBackup(client.ctxWithToken, &common.IDOptions{Id: bid})
	if err != nil || backup == nil {
		client.log.Error().Err(err).Str("backup-id", data.Id()).Msg("Failed to find backup")
		data.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenBackupObject(backup) {
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	data.SetId(backup.GetId())
	return nil
}

// flattenBackupObject creates a map from an Oasis Backup for easy digestion by the Terraform schema.
func flattenBackupObject(backup *backup.Backup) map[string]interface{} {
	return map[string]interface{}{
		backupDataSourceIdFieldName:           backup.GetId(),
		backupDataSourceNameFieldName:         backup.GetName(),
		backupDataSourceDescriptionFieldName:  backup.GetDescription(),
		backupDataSourceURLFieldName:          backup.GetUrl(),
		backupDataSourcePolicyIDFieldName:     backup.GetBackupPolicyId(),
		backupDataSourceDeploymentIDFieldName: backup.GetDeploymentId(),
		backupDataSourceCreatedAtFieldName:    backup.GetCreatedAt().String(),
	}
}
