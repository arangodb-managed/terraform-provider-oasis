//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Joerg Schad, Gergely Brautigam
//

package pkg

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	data "github.com/arangodb-managed/apis/data/v1"
)

const (
	deplOrganizationFieldName              = "organization"
	deplProjectFieldName                   = "project"
	deplNameFieldName                      = "name"
	deplDescriptionFieldName               = "description"
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
	deplConfigurationCoordinatorsFieldName = "coordinators"
	deplConfigurationCoordinatorMemorySize = "coordinator_memory_size"
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

			deplNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},

			deplDescriptionFieldName: {
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
							Required: true,
						},
						// OneShard model
						// Size of nodes being used, e.g., a4
						deplConfigurationNodeSizeIdFieldName: {
							Type:     schema.TypeString,
							Optional: true,
						},
						deplConfigurationNodeCountFieldName: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						deplConfigurationNodeDiskSizeFieldName: {
							Type:     schema.TypeInt,
							Optional: true,
						},

						// Flexible model
						deplConfigurationCoordinatorsFieldName: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						deplConfigurationCoordinatorMemorySize: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						deplConfigurationDbServerCount: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						deplConfigurationDbServerMemorySize: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						deplConfigurationDbServerDiskSize: {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// resourceDeploymentCreate creates an oasis deployment given a project id.
func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	datac := data.NewDataServiceClient(client.conn)
	expandedDepl := expandDeploymentResource(d, client.ProjectID)
	if expandedDepl.Certificates.CaCertificateId == "" {
		cryptoc := crypto.NewCryptoServiceClient(client.conn)
		list, err := cryptoc.ListCACertificates(client.ctxWithToken, &common.ListOptions{ContextId: expandedDepl.GetProjectId()})
		if err != nil {
			client.log.Error().Err(err).Msg("Failed to list CA certificates")
			return err
		}
		if len(list.Items) < 1 {
			client.log.Error().Err(err).Msg("Failed to find any CA certificates")
			return fmt.Errorf("failed to find any CA certificates for project %s", expandedDepl.GetProjectId())
		}
		certificateSelected := false
		for _, c := range list.GetItems() {
			if c.GetIsDefault() {
				expandedDepl.Certificates.CaCertificateId = c.GetId()
				certificateSelected = true
				break
			}
		}
		if !certificateSelected {
			client.log.Error().Err(err).Str("project-id", expandedDepl.ProjectId).Msg("Unable to find default certificate for project. Please select one manually.")
			return fmt.Errorf("Unable to find default certificate for project %s. Please select one manually.", expandedDepl.GetProjectId())
		}
	}

	if depl, err := datac.CreateDeployment(client.ctxWithToken, expandedDepl); err != nil {
		client.log.Error().Err(err).Msg("Failed to create deplyoment.")
		return err
	} else {
		d.SetId(depl.GetId())
	}
	return resourceDeploymentRead(d, m)
}

// location is a convenient wrapper around the location schema for easy parsing
type location struct {
	region   string
	provider string
}

// version is a convenient wrapper around the version schema for easy parsing
type version struct {
	dbVersion     string
	caCertificate string
	ipWhitelist   string
}

// configuration is a convenient wrapper around the configuration schema for easy parsing
type configuration struct {
	model                 string
	nodeSizeId            string
	nodeCount             int
	nodeDiskSize          int
	coordinators          int
	coordinatorMemorySize int
	dbServerCount         int
	dbServerMemorySize    int
	dbServerDiskSize      int
}

// expandDeploymentResource creates an oasis deployment structure out of a terraform schema model.
func expandDeploymentResource(d *schema.ResourceData, defaultProject string) *data.Deployment {
	name := d.Get("name").(string)
	project := defaultProject
	var (
		description string
		ver         version
		loc         location
		conf        configuration
	)
	if v, ok := d.GetOk(deplDescriptionFieldName); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk(deplLocationFieldName); ok {
		loc = expandLocation(v.(*schema.Set))
	}
	if v, ok := d.GetOk(deplVersionFieldName); ok {
		ver = expandVersion(v.(*schema.Set))
	}
	if v, ok := d.GetOk(deplConfigurationFieldName); ok {
		conf = expandConfiguration(v.(*schema.Set))
	}
	var servers *data.Deployment_ServersSpec
	if conf.model == data.ModelFlexible {
		servers = &data.Deployment_ServersSpec{
			Coordinators:          int32(conf.coordinators),
			CoordinatorMemorySize: int32(conf.coordinatorMemorySize),
			Dbservers:             int32(conf.dbServerCount),
			DbserverDiskSize:      int32(conf.dbServerDiskSize),
			DbserverMemorySize:    int32(conf.dbServerMemorySize),
		}
	}

	return &data.Deployment{
		Name:          name,
		Description:   description,
		ProjectId:     project,
		RegionId:      loc.region,
		Version:       ver.dbVersion,
		Certificates:  &data.Deployment_CertificateSpec{CaCertificateId: ver.caCertificate},
		Servers:       servers,
		IpwhitelistId: ver.ipWhitelist,
		Model: &data.Deployment_ModelSpec{
			Model:        conf.model,
			NodeCount:    int32(conf.nodeCount),
			NodeDiskSize: int32(conf.nodeDiskSize),
			NodeSizeId:   conf.nodeSizeId,
		},
	}
}

// expandLocation gathers location data from the location set in terraform schema
func expandLocation(s *schema.Set) (loc location) {
	for _, v := range s.List() {
		item := v.(map[string]interface{})
		if i, ok := item[deplLocationProiderFieldName]; ok {
			loc.provider = i.(string)
		}
		if i, ok := item[deplLocationRegionFieldName]; ok {
			loc.region = i.(string)
		}
	}
	return
}

// expandVersion gathers version data from the version set in terraform schema
func expandVersion(s *schema.Set) (ver version) {
	for _, v := range s.List() {
		item := v.(map[string]interface{})
		if i, ok := item[deplVersionDbVersionFieldName]; ok {
			ver.dbVersion = i.(string)
		}
		if i, ok := item[deplVersionCaCertificateFieldName]; ok {
			ver.caCertificate = i.(string)
		}
		if i, ok := item[deplVersionIpWhitelistFieldName]; ok {
			ver.ipWhitelist = i.(string)
		}
	}
	return
}

// expandConfiguration gathers configuration data from the configuration set in terraform schema
func expandConfiguration(s *schema.Set) (conf configuration) {
	for _, v := range s.List() {
		item := v.(map[string]interface{})
		if i, ok := item[deplConfigurationModelFieldName]; ok {
			conf.model = i.(string)
		}
		if i, ok := item[deplConfigurationNodeSizeIdFieldName]; ok {
			conf.nodeSizeId = i.(string)
		}
		if i, ok := item[deplConfigurationNodeCountFieldName]; ok {
			conf.nodeCount = i.(int)
		}
		if i, ok := item[deplConfigurationNodeDiskSizeFieldName]; ok {
			conf.nodeDiskSize = i.(int)
		}
		if i, ok := item[deplConfigurationCoordinatorsFieldName]; ok {
			conf.coordinators = i.(int)
		}
		if i, ok := item[deplConfigurationCoordinatorMemorySize]; ok {
			conf.coordinators = i.(int)
		}
		if i, ok := item[deplConfigurationDbServerCount]; ok {
			conf.dbServerCount = i.(int)
		}
		if i, ok := item[deplConfigurationDbServerMemorySize]; ok {
			conf.dbServerMemorySize = i.(int)
		}
		if i, ok := item[deplConfigurationDbServerDiskSize]; ok {
			conf.dbServerDiskSize = i.(int)
		}
	}
	return
}

// resourceDeploymentRead retrieves deployment information from terraform stores.
func resourceDeploymentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	datac := data.NewDataServiceClient(client.conn)
	depl, err := datac.GetDeployment(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find deployment")
		d.SetId("")
		return err
	}
	for k, v := range flattenDeployment(depl) {
		if _, ok := d.GetOk(k); ok {
			if err := d.Set(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// flattenDeployment creates a map from a deployment for easy storage on terraform.
func flattenDeployment(depl *data.Deployment) map[string]interface{} {
	conf := flattenConfigurationData(depl)
	loc := flattenLocationData(depl)
	ver := flattenVersionData(depl)

	return map[string]interface{}{
		deplNameFieldName:          depl.GetName(),
		deplProjectFieldName:       depl.GetProjectId(),
		deplConfigurationFieldName: conf,
		deplLocationFieldName:      loc,
		deplVersionFieldName:       ver,
	}
}

// flattenVersionData takes the version part of a deployment and creates a sub map for terraform schema.
func flattenVersionData(depl *data.Deployment) *schema.Set {
	s := &schema.Set{
		F: schema.HashResource(resourceDeployment().Schema[deplVersionFieldName].Elem.(*schema.Resource)),
	}
	versionMap := map[string]interface{}{
		deplVersionDbVersionFieldName:     depl.GetVersion(),
		deplVersionCaCertificateFieldName: depl.GetCertificates().GetCaCertificateId(),
		deplVersionIpWhitelistFieldName:   depl.GetIpwhitelistId(),
	}
	s.Add(versionMap)
	return s
}

// flattenLocationData takes the location part of a deployment and creates a sub map for terraform schema.
func flattenLocationData(depl *data.Deployment) *schema.Set {
	s := &schema.Set{
		F: schema.HashResource(resourceDeployment().Schema[deplLocationFieldName].Elem.(*schema.Resource)),
	}
	locationMap := map[string]interface{}{
		deplLocationRegionFieldName: depl.GetRegionId(),
	}
	s.Add(locationMap)
	return s
}

// flattenConfigurationData takes the configuration part of a deployment and creates a sub map for terraform schema.
func flattenConfigurationData(depl *data.Deployment) *schema.Set {
	s := &schema.Set{
		F: schema.HashResource(resourceDeployment().Schema[deplConfigurationFieldName].Elem.(*schema.Resource)),
	}

	configMap := map[string]interface{}{
		deplConfigurationModelFieldName:        depl.GetModel().GetModel(),
		deplConfigurationCoordinatorsFieldName: int(depl.GetServers().GetCoordinators()),
		deplConfigurationCoordinatorMemorySize: int(depl.GetServers().GetCoordinatorMemorySize()),
		deplConfigurationDbServerCount:         int(depl.GetServers().GetDbservers()),
		deplConfigurationDbServerMemorySize:    int(depl.GetServers().GetDbserverMemorySize()),
		deplConfigurationDbServerDiskSize:      int(depl.GetServers().GetDbserverDiskSize()),
	}
	if depl.GetModel().GetModel() != data.ModelFlexible {
		configMap[deplConfigurationNodeSizeIdFieldName] = depl.GetModel().GetNodeSizeId()
		configMap[deplConfigurationNodeDiskSizeFieldName] = int(depl.GetModel().GetNodeDiskSize())
		configMap[deplConfigurationNodeCountFieldName] = int(depl.GetModel().GetNodeCount())
	}

	s.Add(configMap)
	return s
}

// resourceDeploymentUpdate checks fields for differences and updates a deployment if necessary.
func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	datac := data.NewDataServiceClient(client.conn)
	depl, err := datac.GetDeployment(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to find deployment")
		d.SetId("")
		return err
	}

	if d.HasChange(deplNameFieldName) {
		depl.Name = d.Get(deplNameFieldName).(string)
	}
	if d.HasChange(deplDescriptionFieldName) {
		depl.Description = d.Get(deplDescriptionFieldName).(string)
	}
	if d.HasChange(deplVersionFieldName) {
		// No need to check for individual field change via set.Difference since the values should
		// be that of the deployments if there is no change.
		ver := expandVersion(d.Get(deplVersionFieldName).(*schema.Set))
		depl.Certificates.CaCertificateId = ver.caCertificate
		depl.Version = ver.dbVersion
		depl.IpwhitelistId = ver.ipWhitelist
	}
	if d.HasChange(deplConfigurationFieldName) {
		// No need to check for individual field change via set.Difference since the values should
		// be that of the deployments if there is no change.
		conf := expandConfiguration(d.Get(deplConfigurationFieldName).(*schema.Set))
		depl.Model.Model = conf.model
		depl.Model.NodeSizeId = conf.nodeSizeId
		depl.Model.NodeDiskSize = int32(conf.nodeDiskSize)
		depl.Model.NodeCount = int32(conf.nodeCount)
		depl.Servers.DbserverMemorySize = int32(conf.dbServerMemorySize)
		depl.Servers.DbserverDiskSize = int32(conf.dbServerDiskSize)
		depl.Servers.Dbservers = int32(conf.dbServerCount)
	}

	res, err := datac.UpdateDeployment(client.ctxWithToken, depl)
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to update deployment")
		return err
	}

	d.SetId(res.GetId())
	return resourceDeploymentRead(d, m)
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	datac := data.NewDataServiceClient(client.conn)
	if _, err := datac.DeleteDeployment(client.ctxWithToken, &common.IDOptions{Id: d.Id()}); err != nil {
		client.log.Error().Err(err).Str("deployment-id", d.Id()).Msg("Failed to delete deployment")
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}
