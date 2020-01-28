//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,

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
func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	cryptoc := crypto.NewCryptoServiceClient(client.conn)

	result, err := cryptoc.CreateCACertificate(client.ctxWithToken, expandToCertificate(d))
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create certificate")
		return err
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceCertificateRead(d, m)
}

// resourceCertificateRead handles storing and showing data about the certificate given a stored ID.
// This function should always be called from create and update.
func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	cert, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("certificate-id", d.Id()).Msg("Failed to find certificate")
		return err
	}
	if cert == nil {
		client.log.Error().Str("certificate-id", d.Id()).Msg("Failed to find certificate")
		d.SetId("")
		return nil
	}

	for k, v := range flattenCertificateResource(cert) {
		if _, ok := d.GetOk(k); ok {
			if err := d.Set(k, v); err != nil {
				return err
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
func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	cert, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed get certificate")
		return err
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
		return err
	}
	d.SetId(res.Id)
	return resourceCertificateRead(d, m)
}

// resourceCertificateDelete will be called once the resource is destroyed.
func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	if _, err := cryptoc.DeleteCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("certificate-id", d.Id()).Msg("Failed to delete certificate")
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}
