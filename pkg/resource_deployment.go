//
// DISCLAIMER
//
// Copyright 2020 ArangoDB Inc, Cologne, Germany
//
// Author Joerg Schad, Gergely Brautigam
//

package pkg

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Read:   resourceDeploymentRead,
		Update: resourceDeploymentUpdate,
		Delete: resourceDeploymentDelete,

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{ // If set here, overrides project in provider
				Type:     schema.TypeString,
				Required: true,
			},

			"project": &schema.Schema{ // If set here, overrides project in provider
				Type:     schema.TypeString,
				Optional: true,
			},

			"location": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"version": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ca_certificate": {
							Type:     schema.TypeString,
							Optional: true, // If not set, uses default certificate from project
						},
						"ip_whitelist": {
							Type:     schema.TypeString,
							Optional: true, // If not set, no whitelist is configured
						},
					},
				},
			},

			"configuration": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:     schema.TypeString,
							Required: true,
						},
						// OneShard model
						"node_disk_gb": {
							Type:     schema.TypeInt,
							Required: false,
						},
						// Size of nodes being used, e.g., a4
						"node_size_id": {
							Type:     schema.TypeString,
							Required: false,
						},
						"num_nodes": {
							Type:     schema.TypeString,
							Required: false,
						},

						// Flexible model
						"num_coordinators": {
							Type:     schema.TypeInt,
							Required: false,
						},
						"coordinator_memory_size": {
							Type:     schema.TypeInt,
							Required: false,
						},
						"num_dbservers": {
							Type:     schema.TypeInt,
							Required: false,
						},
						"dbserver_memory_size": {
							Type:     schema.TypeInt,
							Required: false,
						},
						"dbserver_disk_size": {
							Type:     schema.TypeInt,
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
