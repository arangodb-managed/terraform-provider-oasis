//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// AuthorGergely Brautigam
//

package pkg

import (
	"strconv"
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
			"name": { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Required: true,
			},

			"project": { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			"lifetime": { // If set here, overrides project in provider
				Type:     schema.TypeInt,
				Optional: true,
			},

			"well_known_certificate": {
				Type:     schema.TypeBool,
				Optional: true,
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

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	projectId := d.Get("project").(string)
	useWellKnownCertificate := d.Get("well_known_certificate").(bool)
	lifetime, err := strconv.Atoi(d.Get("lifetime").(string))
	if err != nil {
		return err
	}
	var lt *types.Duration
	if lifetime > 0 {
		lt = types.DurationProto(time.Duration(lifetime))
	}
	result, err := cryptoc.CreateCACertificate(client.ctxWithToken, &crypto.CACertificate{
		Name:                    name,
		Description:             description,
		ProjectId:               projectId,
		Lifetime:                lt,
		UseWellKnownCertificate: useWellKnownCertificate,
	})
	if err != nil {
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
		return err
	}
	if cert == nil {
		d.SetId("")
		return nil
	}
	// TODO: Map to schema
	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
