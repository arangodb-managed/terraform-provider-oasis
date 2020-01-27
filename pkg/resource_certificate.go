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

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"project": { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"lifetime": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"use_well_known_certificate": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.Connect()
	if err != nil {
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

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.Connect()
	if err != nil {
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
			err := d.Set(k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// flattenCertificateResource flattens the certificate data into a map interface for easy storage.
func flattenCertificateResource(cert *crypto.CACertificate) map[string]interface{} {
	flatted := map[string]interface{}{
		"name":                       cert.GetName(),
		"description":                cert.GetDescription(),
		"project":                    cert.GetProjectId(),
		"use_well_known_certificate": cert.GetUseWellKnownCertificate(),
		"lifetime":                   int(cert.GetLifetime().GetSeconds()),
		"is_default":                 cert.GetIsDefault(),
		"expires_at":                 cert.GetExpiresAt().String(),
		"created_at":                 cert.GetCreatedAt().String(),
	}
	return flatted
}

// expandToCertificate creates a certificate resource from resource data.
func expandToCertificate(d *schema.ResourceData) *crypto.CACertificate {
	name := d.Get("name").(string)
	projectId := d.Get("project").(string)
	var (
		description             string
		lifetime                int
		useWellKnownCertificate bool
		lt                      *types.Duration
	)
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk("lifetime"); ok {
		lifetime = v.(int)
		if lifetime > 0 {
			lt = types.DurationProto(time.Duration(lifetime))
		}
	}
	if v, ok := d.GetOk("use_well_known_certificate"); ok {
		useWellKnownCertificate = v.(bool)
	}
	return &crypto.CACertificate{
		Name:                    name,
		Description:             description,
		ProjectId:               projectId,
		Lifetime:                lt,
		UseWellKnownCertificate: useWellKnownCertificate,
	}
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.Connect()
	if err != nil {
		return err
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	cert, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		return err
	}
	if cert == nil {
		client.log.Error().Str("certificate-id", d.Id()).Msg("Failed to find certificate")
		d.SetId("")
		return nil
	}

	if d.HasChange("name") {
		cert.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		cert.Description = d.Get("description").(string)
	}
	if d.HasChange("use_well_known_certificate") {
		cert.UseWellKnownCertificate = d.Get("use_well_known_certificate").(bool)
	}
	if d.HasChange("lifetime") {
		cert.Lifetime = types.DurationProto(time.Duration(d.Get("lifetime").(int)))
	}
	res, err := cryptoc.UpdateCACertificate(client.ctxWithToken, cert)
	if err != nil {
		client.log.Error().Err(err).Str("certificate-id", d.Id()).Msg("Failed to update certificate")
		return err
	}
	d.SetId(res.Id)
	return resourceCertificateRead(d, m)
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.Connect()
	if err != nil {
		return err
	}

	cryptoc := crypto.NewCryptoServiceClient(client.conn)
	_, err = cryptoc.DeleteCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("certificate-id", d.Id()).Msg("Failed to delete certificate")
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}
