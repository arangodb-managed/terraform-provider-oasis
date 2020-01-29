//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	security "github.com/arangodb-managed/apis/security/v1"
)

const (
	// IP Whiltelist fields
	ipNameFieldName        = "name"
	ipProjectFieldName     = "project"
	ipDescriptionFieldName = "description"
	ipCIDRRangeFieldName   = "cidr_ranges"
	ipIsDeletedFieldName   = "is_deleted"
	ipCreatedAtFieldName   = "created_at"
)

// resourceIPWhitelist defines the IPWhitelist terraform resource Schema.
func resourceIPWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPWhitelistCreate,
		Read:   resourceIPWhitelistRead,
		Update: resourceIPWhitelistUpdate,
		Delete: resourceIPWhitelistDelete,

		Schema: map[string]*schema.Schema{
			ipNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},

			ipProjectFieldName: { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			ipDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},

			ipCIDRRangeFieldName: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			ipIsDeletedFieldName: {
				Type:     schema.TypeBool,
				Computed: true,
			},

			ipCreatedAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceIPWhitelistCreate handles the creation lifecycle of the IPWhitelist resource
// sets the ID of a given IPWhitelist once the creation is successful. This will be stored in local terraform store.
func resourceIPWhitelistCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	securityc := security.NewSecurityServiceClient(client.conn)
	expanded, err := expandToIPWhitelist(d, client.ProjectID)
	if err != nil {
		return err
	}
	result, err := securityc.CreateIPWhitelist(client.ctxWithToken, expanded)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to create ip whitelist")
		return err
	}
	if result != nil {
		d.SetId(result.Id)
	}
	return resourceIPWhitelistRead(d, m)
}

// expandToIPWhitelist creates an ip whitelist oasis structure out of a terraform schema.
func expandToIPWhitelist(d *schema.ResourceData, defaultProject string) (*security.IPWhitelist, error) {
	var (
		name        string
		description string
		cidrRange   []string
		err         error
	)
	if v, ok := d.GetOk(ipNameFieldName); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk(ipCIDRRangeFieldName); ok {
		cidrRange, err = expandStringList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
	}
	project := defaultProject
	if v, ok := d.GetOk(ipDescriptionFieldName); ok {
		description = v.(string)
	}
	// Overwrite project if it exists
	if v, ok := d.GetOk(ipProjectFieldName); ok {
		project = v.(string)
	}

	return &security.IPWhitelist{
		Name:        name,
		Description: description,
		ProjectId:   project,
		CidrRanges:  cidrRange,
	}, nil
}

// expandStringList creates a string list of items from an interface slice. It also
// verifies if a given string item is empty or not. In case it's empty, an error is thrown.
func expandStringList(list []interface{}) ([]string, error) {
	cidr := make([]string, 0)
	for _, v := range list {
		if v, ok := v.(string); ok {
			if v == "" {
				return []string{}, fmt.Errorf("cidr range cannot be empty")
			}
			cidr = append(cidr, v)
		}
	}
	return cidr, nil
}

// flattenIPWhitelistResource flattens the ip whitelist data into a map interface for easy storage.
func flattenIPWhitelistResource(ip *security.IPWhitelist) map[string]interface{} {
	return map[string]interface{}{
		ipNameFieldName:        ip.GetName(),
		ipProjectFieldName:     ip.GetProjectId(),
		ipDescriptionFieldName: ip.GetDescription(),
		ipCIDRRangeFieldName:   ip.GetCidrRanges(),
		ipCreatedAtFieldName:   ip.GetCreatedAt().String(),
		ipIsDeletedFieldName:   ip.GetIsDeleted(),
	}
}

// resourceIPWhitelistRead handles the read lifecycle of the IPWhitelist resource.
func resourceIPWhitelistRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	securityc := security.NewSecurityServiceClient(client.conn)
	ipWhitelist, err := securityc.GetIPWhitelist(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("ipwhitelist-id", d.Id()).Msg("Failed to find ip whitelist")
		return err
	}
	if ipWhitelist == nil {
		client.log.Error().Str("ipwhitelist-id", d.Id()).Msg("Failed to find ip whitelist")
		d.SetId("")
		return nil
	}

	for k, v := range flattenIPWhitelistResource(ipWhitelist) {
		if _, ok := d.GetOk(k); ok {
			if err := d.Set(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// resourceIPWhitelistDelete will be called once the resource is destroyed.
func resourceIPWhitelistDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	securityc := security.NewSecurityServiceClient(client.conn)
	if _, err := securityc.DeleteIPWhitelist(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("ipwhitelist-id", d.Id()).Msg("Failed to delete ip whitelist")
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

// resourceIPWhitelistUpdate handles the update lifecycle of the IPWhitelist resource.
func resourceIPWhitelistUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	securityc := security.NewSecurityServiceClient(client.conn)
	ipWhitelist, err := securityc.GetIPWhitelist(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Str("ipwhitelist-id", d.Id()).Msg("Failed get ip whitelist")
		return err
	}
	if ipWhitelist == nil {
		client.log.Error().Str("ipwhitelist-id", d.Id()).Msg("Failed to find certificate")
		d.SetId("")
		return nil
	}

	if d.HasChange(ipNameFieldName) {
		ipWhitelist.Name = d.Get(ipNameFieldName).(string)
	}
	if d.HasChange(ipDescriptionFieldName) {
		ipWhitelist.Description = d.Get(ipDescriptionFieldName).(string)
	}
	if d.HasChange(ipCIDRRangeFieldName) {
		cidrRange, err := expandStringList(d.Get(ipCIDRRangeFieldName).([]interface{}))
		if err != nil {
			return err
		}
		ipWhitelist.CidrRanges = cidrRange
	}
	res, err := securityc.UpdateIPWhitelist(client.ctxWithToken, ipWhitelist)
	if err != nil {
		client.log.Error().Err(err).Str("ipwhitelist-id", d.Id()).Msg("Failed to update ip whitelist")
		return err
	}
	d.SetId(res.Id)
	return resourceIPWhitelistRead(d, m)
}
