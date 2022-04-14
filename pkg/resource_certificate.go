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
//

package pkg

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
)

const (
	// Certificate fields
	nameFieldName                    = "name"
	projectFieldName                 = "project"
	descriptionFieldName             = "description"
	lifetimeFieldName                = "lifetime"
	useWellKnownCertificateFieldName = "use_well_known_certificate"
	isDefaultFieldName               = "is_default"
	createdAtFieldName               = "created_at"
	expiresAtFieldName               = "expires_at"
)

// resourceCertificate defines the Certificate terraform resource Schema.
func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		ReadContext:   resourceCertificateRead,
		UpdateContext: resourceCertificateUpdate,
		DeleteContext: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			nameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},

			projectFieldName: { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			descriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},

			lifetimeFieldName: {
				Type:     schema.TypeInt,
				Optional: true,
			},

			useWellKnownCertificateFieldName: {
				Type:     schema.TypeBool,
				Optional: true,
			},

			isDefaultFieldName: {
				Type:     schema.TypeBool,
				Computed: true,
			},

			createdAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			expiresAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceCertificateCreate handles the creation lifecycle of the certificate resource.
// sets the ID of a given certificate once the creation is successful. This will be stored in local terraform store.
func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}
	cryptoc := crypto.NewCryptoServiceClient(client.conn)

	result, err := cryptoc.CreateCACertificate(client.ctxWithToken, expandToCertificate(d))
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create certificate")
		return diag.FromErr(err)
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceCertificateRead(ctx, d, m)
}

// resourceCertificateRead handles storing and showing data about the certificate given a stored ID.
// This function should always be called from create and update.
func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	cert, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("certificate-id", d.Id()).Msg("Failed to find certificate")
		return diag.FromErr(err)
	}
	if cert == nil {
		client.log.Error().Str("certificate-id", d.Id()).Msg("Failed to find certificate")
		d.SetId("")
		return nil
	}

	for k, v := range flattenCertificateResource(cert) {
		if _, ok := d.GetOk(k); ok {
			if err := d.Set(k, v); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

// flattenCertificateResource flattens the certificate data into a map interface for easy storage.
func flattenCertificateResource(cert *crypto.CACertificate) map[string]interface{} {
	flatted := map[string]interface{}{
		nameFieldName:                    cert.GetName(),
		descriptionFieldName:             cert.GetDescription(),
		projectFieldName:                 cert.GetProjectId(),
		useWellKnownCertificateFieldName: cert.GetUseWellKnownCertificate(),
		lifetimeFieldName:                int(cert.GetLifetime().GetSeconds()),
		isDefaultFieldName:               cert.GetIsDefault(),
		expiresAtFieldName:               cert.GetExpiresAt().String(),
		createdAtFieldName:               cert.GetCreatedAt().String(),
	}
	return flatted
}

// expandToCertificate creates a certificate resource from resource data.
func expandToCertificate(d *schema.ResourceData) *crypto.CACertificate {
	n := d.Get(nameFieldName).(string)
	pid := d.Get(projectFieldName).(string)
	var (
		description             string
		lifetime                int
		useWellKnownCertificate bool
		lt                      *types.Duration
	)
	if v, ok := d.GetOk(descriptionFieldName); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk(lifetimeFieldName); ok {
		lifetime = v.(int)
		if lifetime > 0 {
			lt = types.DurationProto(time.Duration(lifetime) * time.Second)
		}
	}
	if v, ok := d.GetOk(useWellKnownCertificateFieldName); ok {
		useWellKnownCertificate = v.(bool)
	}
	return &crypto.CACertificate{
		Name:                    n,
		Description:             description,
		ProjectId:               pid,
		Lifetime:                lt,
		UseWellKnownCertificate: useWellKnownCertificate,
	}
}

// resourceCertificateUpdate handles events in case there is change to the certificate data.
// Only relevant fields are checked for update. Computed fields are ignored.
func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	cert, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed get certificate")
		return diag.FromErr(err)
	}
	if cert == nil {
		client.log.Error().Str("certificate-id", d.Id()).Msg("Failed to find certificate")
		d.SetId("")
		return nil
	}

	if d.HasChange(nameFieldName) {
		cert.Name = d.Get(nameFieldName).(string)
	}
	if d.HasChange(descriptionFieldName) {
		cert.Description = d.Get(descriptionFieldName).(string)
	}
	if d.HasChange(useWellKnownCertificateFieldName) {
		cert.UseWellKnownCertificate = d.Get(useWellKnownCertificateFieldName).(bool)
	}
	if d.HasChange(lifetimeFieldName) {
		cert.Lifetime = types.DurationProto(time.Duration(d.Get(lifetimeFieldName).(int)))
	}
	res, err := cryptoc.UpdateCACertificate(client.ctxWithToken, cert)
	if err != nil {
		client.log.Error().Err(err).Str("certificate-id", d.Id()).Msg("Failed to update certificate")
		return diag.FromErr(err)
	}
	d.SetId(res.Id)
	return resourceCertificateRead(ctx, d, m)
}

// resourceCertificateDelete will be called once the resource is destroyed.
func resourceCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	if _, err := cryptoc.DeleteCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("certificate-id", d.Id()).Msg("Failed to delete certificate")
		return diag.FromErr(err)
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}
