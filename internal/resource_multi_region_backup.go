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
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	backup "github.com/arangodb-managed/apis/backup/v1"
)

const (
	// Multi region Backup field names
	backupRegionIDFieldName       = "region_id"
	backupSourceBackupIDFieldName = "source_backup_id"
)

// resourceMultiRegionBackup defines a Multi Region Backup Oasis resource.
func resourceMultiRegionBackup() *schema.Resource {
	return &schema.Resource{
		Description:   "Oasis Multi Region Backup Resource",
		CreateContext: resourceMultiRegionBackupCreate,
		ReadContext:   resourceBackupRead,
		UpdateContext: resourceBackupUpdate,
		DeleteContext: resourceBackupDelete,
		Schema: map[string]*schema.Schema{
			backupSourceBackupIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Multi Region Backup Resource Region Identifier",
				Optional:    true,
			},
			backupRegionIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Multi Region Backup Resource Region Identifier",
				Optional:    true,
			},

			// backup fields used to return the backup information
			backupNameFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Multi Region Backup Resource Backup Name field, generated based on source backup",
				Computed:    true,
			},
			backupDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Multi Region Backup Resource Backup Description field, generated based on source backup",
				Computed:    true,
			},
			backupUploadFieldName: {
				Type:        schema.TypeBool,
				Description: "Oasis Multi Region Backup Resource Backup Upload field, generated based on source backup",
				Computed:    true,
			},
			backupDeploymentIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Multi Region Backup Resource Backup Deployment ID field, generated based on source backup",
				Computed:    true,
			},
			backupURLFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Multi Region Backup Resource Backup URL field, generated based on source backup",
				Computed:    true,
			},
			backupPolicyIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis Multi Region Backup Resource Backup Policy ID field, generated based on source backup",
				Computed:    true,
			},

			backupAutoDeleteAtFieldName: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Oasis Multi Region Backup Resource Backup Auto Delete At field, generated based on source backup",
			},
		},
	}
}

// resourceMultiRegionBackupCreate will take the schema data from the Terraform config file and call the Oasis client
// to initiate a copy procedure for a given backup and region identifier.
func resourceMultiRegionBackupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	req := &backup.CopyBackupRequest{}
	if v, ok := d.GetOk(backupSourceBackupIDFieldName); ok && strings.TrimSpace(v.(string)) != "" {
		req.SourceBackupId = v.(string)
	} else {
		err := fmt.Errorf("unable to find parse field %s", backupSourceBackupIDFieldName)
		client.log.Error().Err(err).Msg("Source backup identifier required")
		return diag.FromErr(err)
	}
	if v, ok := d.GetOk(backupRegionIDFieldName); ok && strings.TrimSpace(v.(string)) != "" {
		req.RegionId = v.(string)
	} else {
		err := fmt.Errorf("unable to find parse field %s", backupRegionIDFieldName)
		client.log.Error().Err(err).Msg("Region identifier required")
		return diag.FromErr(err)
	}

	backupc := backup.NewBackupServiceClient(client.conn)

	backup, err := backupc.CopyBackup(client.ctxWithToken, req)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create backup")
		return diag.FromErr(err)
	} else {
		d.SetId(backup.GetId())
	}

	for k, v := range flattenBackupResource(backup) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
