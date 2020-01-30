//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Joerg Schad, Gergely Brautigam
//

package pkg

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	data "github.com/arangodb-managed/apis/data/v1"
)

const (
	deplProjectFieldName                         = "project"
	deplNameFieldName                            = "name"
	deplDescriptionFieldName                     = "description"
	deplLocationFieldName                        = "location"
	deplLocationRegionFieldName                  = "region"
	deplVersionAndSecurityFieldName              = "version_and_security"
	deplVersionAndSecurityDbVersionFieldName     = "db_version"
	deplVersionAndSecurityCaCertificateFieldName = "ca_certificate"
	deplVersionAndSecurityIpWhitelistFieldName   = "ip_whitelist"
	deplConfigurationFieldName                   = "configuration"
	deplConfigurationModelFieldName              = "model"
	deplConfigurationNodeSizeIdFieldName         = "node_size_id"
	deplConfigurationNodeCountFieldName          = "node_count"
	deplConfigurationNodeDiskSizeFieldName       = "node_disk_size"
)

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Read:   resourceDeploymentRead,
		Update: resourceDeploymentUpdate,
		Delete: resourceDeploymentDelete,

		Schema: map[string]*schema.Schema{
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
						deplLocationRegionFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			deplVersionAndSecurityFieldName: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplVersionAndSecurityDbVersionFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
						deplVersionAndSecurityCaCertificateFieldName: {
							Type:     schema.TypeString,
							Optional: true, // If not set, uses default certificate from project
						},
						deplVersionAndSecurityIpWhitelistFieldName: {
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

	if len(expandedDepl.Model.NodeSizeId) < 1 && expandedDepl.Model.Model != data.ModelFlexible {
		// Fetch node sizes
		list, err := datac.ListNodeSizes(client.ctxWithToken, &data.NodeSizesRequest{
			ProjectId: expandedDepl.ProjectId,
			RegionId:  expandedDepl.RegionId,
		})
		if err != nil {
			client.log.Fatal().Err(err).Msg("Failed to fetch node size list.")
			return fmt.Errorf("Failed to fetch node size list for region %s", expandedDepl.RegionId)
		}
		if len(list.Items) < 1 {
			client.log.Fatal().Msg("No available node sizes found.")
			return fmt.Errorf("No available node sizes found for region %s", expandedDepl.RegionId)
		}
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].MemorySize < list.Items[j].MemorySize
		})
		expandedDepl.Model.NodeSizeId = list.Items[0].Id
		if expandedDepl.Model.NodeDiskSize == 0 {
			expandedDepl.Model.NodeDiskSize = list.Items[0].MinDiskSize
		}
	}

	if expandedDepl.Model.NodeCount == 0 {
		expandedDepl.Model.NodeCount = 3
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
	region string
}

// version is a convenient wrapper around the version schema for easy parsing
type version struct {
	dbVersion     string
	caCertificate string
	ipWhitelist   string
}

// configuration is a convenient wrapper around the configuration schema for easy parsing
type configuration struct {
	model        string
	nodeSizeId   string
	nodeCount    int
	nodeDiskSize int
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
	if v, ok := d.GetOk(deplVersionAndSecurityFieldName); ok {
		ver = expandVersion(v.(*schema.Set))
	}
	if v, ok := d.GetOk(deplConfigurationFieldName); ok {
		conf = expandConfiguration(v.(*schema.Set))
	}

	return &data.Deployment{
		Name:          name,
		Description:   description,
		ProjectId:     project,
		RegionId:      loc.region,
		Version:       ver.dbVersion,
		Certificates:  &data.Deployment_CertificateSpec{CaCertificateId: ver.caCertificate},
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
		if i, ok := item[deplVersionAndSecurityDbVersionFieldName]; ok {
			ver.dbVersion = i.(string)
		}
		if i, ok := item[deplVersionAndSecurityCaCertificateFieldName]; ok {
			ver.caCertificate = i.(string)
		}
		if i, ok := item[deplVersionAndSecurityIpWhitelistFieldName]; ok {
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
		if i, ok := item[deplConfigurationNodeCountFieldName]; ok && i.(int) != 0 {
			conf.nodeCount = i.(int)
		}
		if i, ok := item[deplConfigurationNodeDiskSizeFieldName]; ok && i.(int) != 0 {
			conf.nodeDiskSize = i.(int)
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

	for k, v := range flattenDeployment(depl, d) {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

// flattenDeployment creates a map from a deployment for easy storage on terraform.
func flattenDeployment(depl *data.Deployment, d *schema.ResourceData) map[string]interface{} {
	conf := flattenConfigurationData(depl, d)
	loc := flattenLocationData(depl)
	ver := flattenVersionData(depl, d)

	return map[string]interface{}{
		deplNameFieldName:               depl.GetName(),
		deplProjectFieldName:            depl.GetProjectId(),
		deplConfigurationFieldName:      conf,
		deplLocationFieldName:           loc,
		deplVersionAndSecurityFieldName: ver,
	}
}

// flattenVersionData takes the version part of a deployment and creates a sub map for terraform schema.
func flattenVersionData(depl *data.Deployment, d *schema.ResourceData) *schema.Set {
	s := &schema.Set{
		F: schema.HashResource(resourceDeployment().Schema[deplVersionAndSecurityFieldName].Elem.(*schema.Resource)),
	}

	versionMap := map[string]interface{}{
		deplVersionAndSecurityDbVersionFieldName:   depl.GetVersion(),
		deplVersionAndSecurityIpWhitelistFieldName: depl.GetIpwhitelistId(),
	}
	// Make sure that certificate is only showing up in the change list and in show
	// if it was explicitly set by the user.
	if v, ok := d.GetOk(deplVersionAndSecurityFieldName); ok {
		m := v.(*schema.Set)
		for _, i := range m.List() {
			inner := i.(map[string]interface{})
			if innerV, ok := inner[deplVersionAndSecurityCaCertificateFieldName]; ok {
				if innerV.(string) != "" {
					if innerV.(string) != depl.GetCertificates().GetCaCertificateId() {
						// Any incoming change will show up as a diff to the user.
						// This is so that we don't miss upstream changes once this field is no longer
						// automatically managed.
						versionMap[deplVersionAndSecurityCaCertificateFieldName] = depl.GetCertificates().GetCaCertificateId()
					} else {
						versionMap[deplVersionAndSecurityCaCertificateFieldName] = innerV.(string)
					}
				}
			}
		}
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
func flattenConfigurationData(depl *data.Deployment, d *schema.ResourceData) *schema.Set {
	s := &schema.Set{
		// Calculate a unique key based on the provided values. This should always calculate the same value
		// for the same changes to avoid conflicting diffs.
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			// All keys added in alphabetical order.
			if v, ok := m[deplConfigurationModelFieldName]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
			if v, ok := m[deplConfigurationNodeCountFieldName]; ok {
				buf.WriteString(fmt.Sprintf("%d-", v.(int)))
			}
			if v, ok := m[deplConfigurationNodeDiskSizeFieldName]; ok {
				buf.WriteString(fmt.Sprintf("%d-", v.(int)))
			}
			if v, ok := m[deplConfigurationNodeSizeIdFieldName]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
			return hashcode.String(buf.String())
		},
	}

	configMap := map[string]interface{}{
		deplConfigurationModelFieldName:        depl.GetModel().GetModel(),
		deplConfigurationNodeSizeIdFieldName:   depl.GetModel().GetNodeSizeId(),
		deplConfigurationNodeDiskSizeFieldName: int(depl.GetModel().GetNodeDiskSize()),
		deplConfigurationNodeCountFieldName:    int(depl.GetModel().GetNodeCount()),
	}

	// Remove the keys that the user did not set up manually. But always check if there is
	// a potential upstream change.
	if v, ok := d.GetOk(deplConfigurationFieldName); ok {
		m := v.(*schema.Set)
		for _, v := range m.List() {
			change := v.(map[string]interface{})
			if v, ok := change[deplConfigurationNodeCountFieldName]; ok {
				if v.(int) != 0 {
					if v.(int) == int(depl.GetModel().GetNodeCount()) {
						configMap[deplConfigurationNodeCountFieldName] = v.(int)
					}
				} else {
					delete(configMap, deplConfigurationNodeCountFieldName)
				}
			}
			if v, ok := change[deplConfigurationNodeSizeIdFieldName]; ok {
				if v.(string) != "" {
					if v.(string) == depl.GetModel().GetNodeSizeId() {
						configMap[deplConfigurationNodeSizeIdFieldName] = v.(string)
					}
				} else {
					delete(configMap, deplConfigurationNodeSizeIdFieldName)
				}
			}
			if v, ok := change[deplConfigurationNodeDiskSizeFieldName]; ok {
				if v.(int) != 0 {
					if v.(int) == int(depl.GetModel().GetNodeDiskSize()) {
						configMap[deplConfigurationNodeDiskSizeFieldName] = v.(int)
					}
				} else {
					delete(configMap, deplConfigurationNodeDiskSizeFieldName)
				}
			}
		}
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
	if d.HasChange(deplVersionAndSecurityFieldName) {
		ver := expandVersion(d.Get(deplVersionAndSecurityFieldName).(*schema.Set))
		if ver.caCertificate != "" {
			depl.Certificates.CaCertificateId = ver.caCertificate
		}
		if ver.dbVersion != "" {
			depl.Version = ver.dbVersion
		}
		if ver.ipWhitelist != "" {
			depl.IpwhitelistId = ver.ipWhitelist
		}
	}
	if d.HasChange(deplConfigurationFieldName) {
		conf := expandConfiguration(d.Get(deplConfigurationFieldName).(*schema.Set))

		if conf.model != "" {
			depl.Model.Model = conf.model
		}
		if conf.nodeSizeId != "" {
			depl.Model.NodeSizeId = conf.nodeSizeId
		}
		if conf.nodeDiskSize != 0 {
			depl.Model.NodeDiskSize = int32(conf.nodeDiskSize)
		}
		if conf.nodeCount != 0 {
			depl.Model.NodeCount = int32(conf.nodeCount)
		}
	}

	if res, err := datac.UpdateDeployment(client.ctxWithToken, depl); err != nil {
		client.log.Error().Err(err).Msg("Failed to update deployment")
		return err
	} else {
		d.SetId(res.GetId())
	}

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
