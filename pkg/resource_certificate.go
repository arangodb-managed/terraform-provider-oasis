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
	name                    = "name"
	project                 = "project"
	description             = "description"
	lifetime                = "lifetime"
	useWellKnownCertificate = "use_well_known_certificate"
	isDefault               = "is_default"
	createdAt               = "created_at"
	expiresAt               = "expires_at"
)

// resourceCertificate defines the Certificate terraform resource Schema.
func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			name: {
				Type:     schema.TypeString,
				Required: true,
			},

			project: { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			description: {
				Type:     schema.TypeString,
				Optional: true,
			},

			lifetime: {
				Type:     schema.TypeInt,
				Optional: true,
			},

			useWellKnownCertificate: {
				Type:     schema.TypeBool,
				Optional: true,
			},

			isDefault: {
				Type:     schema.TypeBool,
				Computed: true,
			},

			createdAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			expiresAt: {
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
		name:                    cert.GetName(),
		description:             cert.GetDescription(),
		project:                 cert.GetProjectId(),
		useWellKnownCertificate: cert.GetUseWellKnownCertificate(),
		lifetime:                int(cert.GetLifetime().GetSeconds()),
		isDefault:               cert.GetIsDefault(),
		expiresAt:               cert.GetExpiresAt().String(),
		createdAt:               cert.GetCreatedAt().String(),
	}
	return flatted
}

// expandToCertificate creates a certificate resource from resource data.
func expandToCertificate(d *schema.ResourceData) *crypto.CACertificate {
	n := d.Get(name).(string)
	pid := d.Get(project).(string)
	var (
		desc         string
		lifeTime     int
		useWellKnown bool
		lt           *types.Duration
	)
	if v, ok := d.GetOk(description); ok {
		desc = v.(string)
	}
	if v, ok := d.GetOk(lifetime); ok {
		lifeTime = v.(int)
		if lifeTime > 0 {
			lt = types.DurationProto(time.Duration(lifeTime) * time.Second)
		}
	}
	if v, ok := d.GetOk(useWellKnownCertificate); ok {
		useWellKnown = v.(bool)
	}
	return &crypto.CACertificate{
		Name:                    n,
		Description:             desc,
		ProjectId:               pid,
		Lifetime:                lt,
		UseWellKnownCertificate: useWellKnown,
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

	if d.HasChange(name) {
		cert.Name = d.Get(name).(string)
	}
	if d.HasChange(description) {
		cert.Description = d.Get(description).(string)
	}
	if d.HasChange(useWellKnownCertificate) {
		cert.UseWellKnownCertificate = d.Get(useWellKnownCertificate).(bool)
	}
	if d.HasChange(lifetime) {
		cert.Lifetime = types.DurationProto(time.Duration(d.Get(lifetime).(int)))
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
