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
)

const (
	backupSourceBackupIDFieldName = "source_backup_id"
	backupRegionIDFieldName       = "region_id"
)

// resourceBackup defines a Multi Region Backup Oasis resource.
func resourceMultiRegionBackup() *schema.Resource {
	return &schema.Resource{
		Description:   "Oasis Backup Resource",
		CreateContext: resourceCopyBackup,
		Schema: map[string]*schema.Schema{
			backupSourceBackupIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis backup resource identifier field",
				Required:    true,
			},
			backupRegionIDFieldName: {
				Type:        schema.TypeString,
				Description: "Oasis cloud provider region identifier field",
				Required:    true,
			},
		},
	}
}

// resourceCopyBackup will copy backup resource to given region.
func resourceCopyBackup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	backupc := backup.NewBackupServiceClient(client.conn)
	sourceBackupID := ""
	if v, ok := d.GetOk(backupSourceBackupIDFieldName); ok {
		sourceBackupID = v.(string)
	} else {
		return diag.Errorf("Source backup identifier required")
	}

	regionID := ""
	if v, ok := d.GetOk(backupRegionIDFieldName); ok {
		regionID = v.(string)
	} else {
		return diag.Errorf("Region identifier required")
	}

	bu, err := backupc.CopyBackup(client.ctxWithToken, &backup.CopyBackupRequest{
		SourceBackupId: sourceBackupID,
		RegionId:       regionID,
	})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to copy backup")
		d.SetId("")
		return diag.FromErr(err)
	} else {
		d.SetId(bu.GetId())
	}
	return resourceBackupRead(ctx, d, m)
}
