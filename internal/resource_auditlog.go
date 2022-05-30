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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	audit "github.com/arangodb-managed/apis/audit/v1"
	common "github.com/arangodb-managed/apis/common/v1"
)

const (
	// Audit Log field names
	auditLogNameFieldName         = "name"
	auditLogDescriptionFieldName  = "description"
	auditLogOrganizationFieldName = "organization"
	auditLogIsDefaultFieldName    = "is_default"
)

// resourceAuditLog defines an Oasis Audit Log resource.
func resourceAuditLog() *schema.Resource {
	return &schema.Resource{
		Description:   "Oasis Audit Log Resource",
		CreateContext: resourceAuditLogCreate,
		ReadContext:   resourceAuditLogRead,
		UpdateContext: resourceAuditLogUpdate,
		DeleteContext: resourceAuditLogDelete,
		Schema: map[string]*schema.Schema{
			auditLogNameFieldName: {
				Type:        schema.TypeString,
				Description: "Audit Log Resource Audit Log Name field",
				Required:    true,
			},
			auditLogDescriptionFieldName: {
				Type:        schema.TypeString,
				Description: "Audit Log Resource Audit Log Description field",
				Optional:    true,
			},
			auditLogOrganizationFieldName: {
				Type:        schema.TypeString,
				Description: "Audit Log Resource Organization ID field",
				Optional:    true,
			},
			auditLogIsDefaultFieldName: {
				Type:        schema.TypeBool,
				Description: "Audit Log Resource Audit Log Is Default field",
				Optional:    true,
			},
		},
	}
}

// resourceAuditLogRead will gather information from the Terraform store for Oasis Audit Log resource and display it accordingly.
func resourceAuditLogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	auditc := audit.NewAuditServiceClient(client.conn)
	auditLog, err := auditc.GetAuditLog(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil || auditLog == nil {
		client.log.Error().Err(err).Str("auditLog-id", d.Id()).Msg("Failed to find AuditLog")
		d.SetId("")
		return diag.FromErr(err)
	}

	for k, v := range flattenAuditLogResource(auditLog) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// resourceAuditLogCreate will take the schema data from the Terraform config file and call the Oasis client
// to initiate a create procedure for an Audit Log. It will call helper methods to construct the necessary data
// in order to create this object.
func resourceAuditLogCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	auditc := audit.NewAuditServiceClient(client.conn)
	expanded, err := expandAuditLogResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := auditc.CreateAuditLog(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create auditlog")
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk(auditLogIsDefaultFieldName); ok {
		isDefault := v.(bool)
		if isDefault {
			if _, err := auditc.SetDefaultAuditLog(client.ctxWithToken, &audit.SetDefaultAuditLogRequest{
				OrganizationId: result.GetOrganizationId(),
				AuditlogId:     result.GetId(),
			}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if result != nil {
		d.SetId(result.Id)
	}
	return resourceAuditLogRead(ctx, d, m)
}

// expandAuditLogResource will take a Terraform flat map schema data and turn it into an Oasis Audit Log.
func expandAuditLogResource(d *schema.ResourceData) (*audit.AuditLog, error) {
	ret := &audit.AuditLog{}
	if v, ok := d.GetOk(auditLogNameFieldName); ok {
		ret.Name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", auditLogNameFieldName)
	}
	if v, ok := d.GetOk(auditLogDescriptionFieldName); ok {
		ret.Description = v.(string)
	}
	if v, ok := d.GetOk(auditLogOrganizationFieldName); ok {
		ret.OrganizationId = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", auditLogOrganizationFieldName)
	}

	ret.Destinations = []*audit.AuditLog_Destination{{Type: audit.DestinationCloud}}

	return ret, nil
}

// resourceAuditLogDelete will delete a given AuditLog resource based on the given ID
func resourceAuditLogDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	auditc := audit.NewAuditServiceClient(client.conn)
	if _, err := auditc.DeleteAuditLog(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("auditlog-id", d.Id()).Msg("Failed to delete Audit Log")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceAuditLogUpdate will take a resource diff and apply changes accordingly if there are any.
func resourceAuditLogUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	auditc := audit.NewAuditServiceClient(client.conn)
	auditLog, err := auditc.GetAuditLog(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find AuditLog")
		d.SetId("")
		return diag.FromErr(err)
	}
	isChangedDefault := false
	newDefaultAuditLogID := ""
	// Main fields
	if d.HasChange(auditLogNameFieldName) {
		auditLog.Name = d.Get(auditLogNameFieldName).(string)
	}
	if d.HasChange(auditLogDescriptionFieldName) {
		auditLog.Description = d.Get(auditLogDescriptionFieldName).(string)
	}
	if d.HasChange(auditLogOrganizationFieldName) {
		return diag.FromErr(errors.New("organization id cannot be changed"))
	}
	if d.HasChange(auditLogIsDefaultFieldName) {
		isDefaultAuditLog := d.Get(auditLogIsDefaultFieldName).(bool)
		if auditLog.GetIsDefault() != isDefaultAuditLog {
			isChangedDefault = true
			if isDefaultAuditLog {
				newDefaultAuditLogID = auditLog.GetId()
			} else {
				newDefaultAuditLogID = ""
			}
		}
	}

	res, err := auditc.UpdateAuditLog(client.ctxWithToken, auditLog)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update AuditLog")
		return diag.FromErr(err)
	}

	if isChangedDefault {
		if _, err := auditc.SetDefaultAuditLog(client.ctxWithToken, &audit.SetDefaultAuditLogRequest{
			OrganizationId: auditLog.GetOrganizationId(),
			AuditlogId:     newDefaultAuditLogID,
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(res.GetId())

	return resourceAuditLogRead(ctx, d, m)
}

// flattenAuditLogResource will take an AuditLog object and turn it into a flat map for terraform digestion.
func flattenAuditLogResource(auditLog *audit.AuditLog) map[string]interface{} {
	return map[string]interface{}{
		auditLogNameFieldName:         auditLog.GetName(),
		auditLogDescriptionFieldName:  auditLog.GetDescription(),
		auditLogOrganizationFieldName: auditLog.GetOrganizationId(),
	}
}
