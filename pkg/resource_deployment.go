//
// DISCLAIMER
//
// Copyright 2020-2021 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
// Author Robert Stam
// Author Marcin Swiderski
//

package pkg

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	deplTAndCAcceptedFieldName             = "terms_and_conditions_accepted"
	deplProjectFieldName                   = "project"
	deplNameFieldName                      = "name"
	deplDescriptionFieldName               = "description"
	deplLocationFieldName                  = "location"
	deplLocationRegionFieldName            = "region"
	deplVersionFieldName                   = "version"
	deplVersionDbVersionFieldName          = "db_version"
	deplSecurityFieldName                  = "security"
	deplSecurityCaCertificateFieldName     = "ca_certificate"
	deplSecurityIpAllowlistFieldName       = "ip_allowlist"
	deplConfigurationFieldName             = "configuration"
	deplConfigurationModelFieldName        = "model"
	deplConfigurationNodeSizeIdFieldName   = "node_size_id"
	deplConfigurationNodeCountFieldName    = "node_count"
	deplConfigurationNodeDiskSizeFieldName = "node_disk_size"
)

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Read:   resourceDeploymentRead,
		Update: resourceDeploymentUpdate,
		Delete: resourceDeploymentDelete,

		Schema: map[string]*schema.Schema{
			deplTAndCAcceptedFieldName: {
				Type:     schema.TypeBool,
				Required: true,
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
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplLocationRegionFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			deplVersionFieldName: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "1" && new == "0"
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplVersionDbVersionFieldName: {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
						},
					},
				},
			},

			deplSecurityFieldName: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "1" && new == "0"
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplSecurityCaCertificateFieldName: {
							Type:     schema.TypeString,
							Optional: true, // If not set, uses default certificate from project
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
						},
						deplSecurityIpAllowlistFieldName: {
							Type:     schema.TypeString,
							Optional: true, // If not set, no allowlist is configured
						},
					},
				},
			},

			deplConfigurationFieldName: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deplConfigurationModelFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
						deplConfigurationNodeSizeIdFieldName: {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
						},
						deplConfigurationNodeCountFieldName: {
							Type:     schema.TypeInt,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == "0"
							},
						},
						deplConfigurationNodeDiskSizeFieldName: {
							Type:     schema.TypeInt,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == "0"
							},
						},
					},
				},
			},
		},
	}
}

// resourceDeploymentCreate creates an oasis deployment given a project id.
// It will automatically select a certificate if none is provided and will
// automatically select the smallest node size if none is provided.
func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	// Check if the T&C has been accepted
	if v, ok := d.GetOk(deplTAndCAcceptedFieldName); ok {
		if !v.(bool) {
			client.log.Error().Str("name", deplTAndCAcceptedFieldName).Msg("Field should be set to accept Terms and Conditions")
			return fmt.Errorf("Field '%s' should be set to accept Terms and Conditions", deplTAndCAcceptedFieldName)
		}
	} else {
		client.log.Error().Str("name", deplTAndCAcceptedFieldName).Msg("Unable to find field, which is required to accept Terms and Conditions")
		return fmt.Errorf("Unable to find field %s", deplTAndCAcceptedFieldName)
	}

	datac := data.NewDataServiceClient(client.conn)
	expandedDepl, err := expandDeploymentResource(d, client.ProjectID)
	if err != nil {
		return err
	}
	if expandedDepl.Version == "" {
		defaultVersion, err := datac.GetDefaultVersion(client.ctxWithToken, &common.Empty{})
		if err != nil {
			client.log.Error().Err(err).Msg("Failed to get default version")
			return err
		}
		expandedDepl.Version = defaultVersion.Version
	}
	if expandedDepl.Certificates.CaCertificateId == "" {
		cryptoc := crypto.NewCryptoServiceClient(client.conn)
		list, err := cryptoc.ListCACertificates(client.ctxWithToken, &common.ListOptions{ContextId: expandedDepl.GetProjectId()})
		if err != nil {
			client.log.Error().Err(err).Msg("Failed to list CA certificates")
			return err
		}
		if len(list.GetItems()) < 1 {
			client.log.Error().Err(err).Msg("Failed to find any CA certificates")
			return fmt.Errorf("failed to find any CA certificates for project %s", expandedDepl.GetProjectId())
		}
		// Select the default certificate
		for _, c := range list.GetItems() {
			if c.GetIsDefault() {
				expandedDepl.Certificates.CaCertificateId = c.GetId()
				break
			}
		}

		// If the list is one item long, select it, regardless of it being the default or not.
		if len(list.GetItems()) == 1 {
			expandedDepl.Certificates.CaCertificateId = list.GetItems()[0].GetId()
		}

		if expandedDepl.Certificates.CaCertificateId == "" {
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
	if expandedDepl.GetModel().GetModel() == data.ModelDeveloper {
		expandedDepl.Model.NodeCount = 1
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	proj, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: expandedDepl.GetProjectId()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to get project")
		return err
	}
	tAndC, err := rmc.GetCurrentTermsAndConditions(client.ctxWithToken, &common.IDOptions{Id: proj.GetOrganizationId()})
	if err != nil {
		client.log.Error().Err(err).Msg("Failed to get Terms and Conditions")
		return err
	}
	client.log.Info().Str("id", tAndC.GetId()).Msg("Terms and Conditions are accepted")
	expandedDepl.AcceptedTermsAndConditionsId = tAndC.GetId()

	if depl, err := datac.CreateDeployment(client.ctxWithToken, expandedDepl); err != nil {
		client.log.Error().Err(err).Msg("Failed to create deployment.")
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
	dbVersion string
}

// security is a convenient wrapper around the security schema for easy parsing
type securityFields struct {
	caCertificate string
	ipAllowlist   string
}

// configuration is a convenient wrapper around the configuration schema for easy parsing
type configuration struct {
	model        string
	nodeSizeId   string
	nodeCount    int
	nodeDiskSize int
}

// expandDeploymentResource creates an oasis deployment structure out of a terraform schema model.
func expandDeploymentResource(d *schema.ResourceData, defaultProject string) (*data.Deployment, error) {
	project := defaultProject
	var (
		name        string
		description string
		ver         version
		loc         location
		conf        configuration
		sec         securityFields
		err         error
	)
	if v, ok := d.GetOk(deplNameFieldName); ok {
		name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", deplNameFieldName)
	}
	if v, ok := d.GetOk(deplDescriptionFieldName); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk(deplProjectFieldName); ok {
		project = v.(string)
	}
	if v, ok := d.GetOk(deplLocationFieldName); ok {
		if loc, err = expandLocation(v.([]interface{})); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", deplLocationFieldName)
	}
	if v, ok := d.GetOk(deplVersionFieldName); ok {
		if ver, err = expandVersion(v.([]interface{})); err != nil {
			return nil, err
		}
	}
	if v, ok := d.GetOk(deplSecurityFieldName); ok {
		sec = expandSecurity(v.([]interface{}))
	}
	if v, ok := d.GetOk(deplConfigurationFieldName); ok {
		if conf, err = expandConfiguration(v.([]interface{})); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", deplConfigurationFieldName)
	}

	return &data.Deployment{
		Name:          name,
		Description:   description,
		ProjectId:     project,
		RegionId:      loc.region,
		Version:       ver.dbVersion,
		Certificates:  &data.Deployment_CertificateSpec{CaCertificateId: sec.caCertificate},
		IpallowlistId: sec.ipAllowlist,
		Model: &data.Deployment_ModelSpec{
			Model:        conf.model,
			NodeCount:    int32(conf.nodeCount),
			NodeDiskSize: int32(conf.nodeDiskSize),
			NodeSizeId:   conf.nodeSizeId,
		},
	}, nil
}

// expandLocation gathers location data from the terraform store
func expandLocation(s []interface{}) (loc location, err error) {
	for _, v := range s {
		if item, ok := v.(map[string]interface{}); ok {
			if i, ok := item[deplLocationRegionFieldName]; ok {
				loc.region = i.(string)
			} else {
				return loc, fmt.Errorf("failed to parse field %s", deplLocationFieldName)
			}
		}
	}
	return
}

// expandVersion gathers version and security data from the terraform store
func expandVersion(s []interface{}) (ver version, err error) {
	for _, v := range s {
		if item, ok := v.(map[string]interface{}); ok {
			if i, ok := item[deplVersionDbVersionFieldName]; ok {
				ver.dbVersion = i.(string)
			}
		}
	}
	return
}

// expandSecurity gathers security data from the terraform store
func expandSecurity(s []interface{}) (sec securityFields) {
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[deplSecurityCaCertificateFieldName]; ok {
			sec.caCertificate = i.(string)
		}
		if i, ok := item[deplSecurityIpAllowlistFieldName]; ok {
			sec.ipAllowlist = i.(string)
		}
	}
	return
}

// expandConfiguration gathers configuration data from the configuration set in terraform schema
func expandConfiguration(s []interface{}) (conf configuration, err error) {
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[deplConfigurationModelFieldName]; ok {
			conf.model = i.(string)
		} else {
			return conf, fmt.Errorf("failed to parse field %s", deplConfigurationModelFieldName)
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

	for k, v := range flattenDeployment(depl) {
		if err := d.Set(k, v); err != nil {
			d.SetId("")
			return err
		}
	}
	return nil
}

// flattenDeployment creates a map from a deployment for easy storage on terraform.
func flattenDeployment(depl *data.Deployment) map[string]interface{} {
	conf := flattenConfigurationData(depl)
	loc := flattenLocationData(depl)
	ver := flattenVersion(depl)
	sec := flattenSecurity(depl)

	return map[string]interface{}{
		deplNameFieldName:          depl.GetName(),
		deplProjectFieldName:       depl.GetProjectId(),
		deplDescriptionFieldName:   depl.GetDescription(),
		deplConfigurationFieldName: conf,
		deplLocationFieldName:      loc,
		deplVersionFieldName:       ver,
		deplSecurityFieldName:      sec,
	}
}

// flattenVersion takes the version part of a deployment and creates a sub map for terraform schema.
func flattenVersion(depl *data.Deployment) []interface{} {
	return []interface{}{
		map[string]interface{}{
			deplVersionDbVersionFieldName: depl.GetVersion(),
		},
	}
}

// flattenSecurity takes the security part of a deployment and creates a sub map for terraform schema.
func flattenSecurity(depl *data.Deployment) []interface{} {
	return []interface{}{
		map[string]interface{}{
			deplSecurityIpAllowlistFieldName:   depl.GetIpallowlistId(),
			deplSecurityCaCertificateFieldName: depl.GetCertificates().GetCaCertificateId(),
		},
	}
}

// flattenLocationData takes the location part of a deployment and creates a sub map for terraform schema.
func flattenLocationData(depl *data.Deployment) []interface{} {
	return []interface{}{
		map[string]interface{}{
			deplLocationRegionFieldName: depl.GetRegionId(),
		},
	}
}

// flattenConfigurationData takes the configuration part of a deployment and creates a sub map for terraform schema.
func flattenConfigurationData(depl *data.Deployment) []interface{} {
	return []interface{}{
		map[string]interface{}{
			deplConfigurationModelFieldName:        depl.GetModel().GetModel(),
			deplConfigurationNodeSizeIdFieldName:   depl.GetModel().GetNodeSizeId(),
			deplConfigurationNodeDiskSizeFieldName: int(depl.GetModel().GetNodeDiskSize()),
			deplConfigurationNodeCountFieldName:    int(depl.GetModel().GetNodeCount()),
		},
	}
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
		ver, err := expandVersion(d.Get(deplVersionFieldName).([]interface{}))
		if err != nil {
			return err
		}
		if ver.dbVersion != "" {
			depl.Version = ver.dbVersion
		}
	}
	if d.HasChange(deplSecurityFieldName) {
		sec := expandSecurity(d.Get(deplSecurityFieldName).([]interface{}))
		if sec.caCertificate != "" {
			depl.Certificates.CaCertificateId = sec.caCertificate
		}
		if sec.ipAllowlist != "" {
			depl.IpallowlistId = sec.ipAllowlist
		}
	}
	if d.HasChange(deplConfigurationFieldName) {
		conf, err := expandConfiguration(d.Get(deplConfigurationFieldName).([]interface{}))
		if err != nil {
			return err
		}

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
