//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Joerg Schad, Gergely Brautigam
//

package pkg

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	data "github.com/arangodb-managed/apis/data/v1"
)

const (
	deplOrganizationFieldName              = "organization"
	deplProjectFieldName                   = "project"
	deplLocationFieldName                  = "location"
	deplLocationProiderFieldName           = "provider"
	deplLocationRegionFieldName            = "region"
	deplVersionFieldName                   = "version"
	deplVersionDbVersionFieldName          = "db_version"
	deplVersionCaCertificateFieldName      = "ca_certificate"
	deplVersionIpWhitelistFieldName        = "ip_whitelist"
	deplConfigurationFieldName             = "configuration"
	deplConfigurationModelFieldName        = "model"
	deplConfigurationNodeSizeIdFieldName   = "node_size_id"
	deplConfigurationNodeCountFieldName    = "node_count"
	deplConfigurationNodeDiskSizeFieldName = "node_disk_size"
	deplConfigurationCoordinatrosFieldName = "coordinators"
	deplConfigurationCoordinatorMemroySize = "coordinator_memory_size"
	deplConfigurationDbServerCount         = "dbserver_count"
	deplConfigurationDbServerMemorySize    = "dbserver_memory_size"
	deplConfigurationDbServerDiskSize      = "dbserver_disk_size"
)

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Read:   resourceDeploymentRead,
		Update: resourceDeploymentUpdate,
		Delete: resourceDeploymentDelete,

		Schema: map[string]*schema.Schema{
			deplOrganizationFieldName: { // If set here, overrides orgranization in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			deplProjectFieldName: { // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			deplLocationFieldName: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplLocationProiderFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
						deplLocationRegionFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			deplVersionFieldName: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplVersionDbVersionFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
						deplVersionCaCertificateFieldName: {
							Type:     schema.TypeString,
							Optional: true, // If not set, uses default certificate from project
						},
						deplVersionIpWhitelistFieldName: {
							Type:     schema.TypeString,
							Optional: true, // If not set, no whitelist is configured
						},
					},
				},
			},

			deplConfigurationFieldName: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplConfigurationModelFieldName: {
							Type:     schema.TypeString,
							Default:  data.ModelOneShard,
							Required: true,
						},
						// OneShard model
						// Size of nodes being used, e.g., a4
						deplConfigurationNodeSizeIdFieldName: {
							Type:     schema.TypeString,
							Required: false,
						},
						deplConfigurationNodeCountFieldName: {
							Type:     schema.TypeInt,
							Default:  3,
							Required: false,
						},
						deplConfigurationNodeDiskSizeFieldName: {
							Type:     schema.TypeInt,
							Default:  0,
							Required: false,
						},

						// Flexible model
						deplConfigurationCoordinatrosFieldName: {
							Type:     schema.TypeInt,
							Default:  3,
							Required: false,
						},
						deplConfigurationCoordinatorMemroySize: {
							Type:     schema.TypeInt,
							Default:  4,
							Required: false,
						},
						deplConfigurationDbServerCount: {
							Type:     schema.TypeInt,
							Default:  3,
							Required: false,
						},
						deplConfigurationDbServerMemorySize: {
							Type:     schema.TypeInt,
							Default:  4,
							Required: false,
						},
						deplConfigurationDbServerDiskSize: {
							Type:     schema.TypeInt,
							Default:  32,
							Required: false,
						},
					},
				},
			},
		},
	}
}

func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	err := m.(*Client).Connect()
	if err != nil {
		return err
	}
	return resourceDeploymentRead(d, m)
}

func resourceDeploymentRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
