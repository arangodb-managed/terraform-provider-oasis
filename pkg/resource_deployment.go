//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Joerg Schad, Gergely Brautigam
//

package pkg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
)

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

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
	client := m.(*Client)

	err := client.Connect()
	if err != nil {
		return err
	}

	datac := data.NewDataServiceClient(client.conn)
	deployment, err := datac.CreateDeployment(client.ctxWithToken, &data.Deployment{
		ProjectId:   d.Get("project").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		RegionId:    d.Get("location").(*schema.ResourceData).Get("region").(string),
		Version:     d.Get("version").(*schema.ResourceData).Get("db_version").(string),
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: d.Get("version").(*schema.ResourceData).Get("ca_certificate").(string),
		},
		IpwhitelistId: d.Get("version").(*schema.ResourceData).Get("ip_whitelist").(string),
		Servers: &data.Deployment_ServersSpec{
			Coordinators:          d.Get("configuration").(*schema.ResourceData).Get("num_coordinators").(int32),
			CoordinatorMemorySize: d.Get("configuration").(*schema.ResourceData).Get("coordinator_memory_size").(int32),
			Dbservers:             d.Get("configuration").(*schema.ResourceData).Get("num_dbservers").(int32),
			DbserverMemorySize:    d.Get("configuration").(*schema.ResourceData).Get("dbserver_memory_size").(int32),
			DbserverDiskSize:      d.Get("configuration").(*schema.ResourceData).Get("dbserver_disk_size").(int32),
		},
		Model: &data.Deployment_ModelSpec{
			Model:        d.Get("configuration").(*schema.ResourceData).Get("model").(string),
			NodeSizeId:   d.Get("configuration").(*schema.ResourceData).Get("node_size_id").(string),
			NodeCount:    d.Get("configuration").(*schema.ResourceData).Get("num_nodes").(int32),
			NodeDiskSize: d.Get("configuration").(*schema.ResourceData).Get("node_disk_gb").(int32),
		},
	})
	if err != nil {
		return err
	}

	d.SetId(deployment.Id)
	return resourceDeploymentRead(d, m)
}

func resourceDeploymentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	err := client.Connect()
	if err != nil {
		return err
	}

	datac := data.NewDataServiceClient(client.conn)

	deployment, err := datac.GetDeployment(context.Background(), &common.IDOptions{Id: d.Id()})

	if err != nil {
		return err
	}
	if deployment == nil {
		d.SetId("")
		return nil
	}
	// TODO: Map schema to deployment here.
	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
