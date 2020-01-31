package pkg

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/stretchr/testify/assert"

	data "github.com/arangodb-managed/apis/data/v1"
)

func TestExpandingDeploymentResource(t *testing.T) {
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "123456789",
		RegionId:    "gcp-europe-west4",
		Version:     "3.6.0",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpwhitelistId: "ip-whitelist",
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
	}
	flattened := flattenDeployment(depl)
	expected := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionAndSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplVersionAndSecurityDbVersionFieldName:     "3.6.0",
				deplVersionAndSecurityCaCertificateFieldName: "certificate-id",
				deplVersionAndSecurityIpWhitelistFieldName:   "ip-whitelist",
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:        "oneshard",
				deplConfigurationNodeSizeIdFieldName:   "a8",
				deplConfigurationNodeCountFieldName:    3,
				deplConfigurationNodeDiskSizeFieldName: 32,
			},
		},
	}
	assert.Equal(t, expected, flattened)
}

func TestExpandDeployment(t *testing.T) {
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "123456789",
		RegionId:    "gcp-europe-west4",
		Version:     "3.6.0",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpwhitelistId: "ip-whitelist",
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
	}
	raw := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionAndSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplVersionAndSecurityDbVersionFieldName:     "3.6.0",
				deplVersionAndSecurityCaCertificateFieldName: "certificate-id",
				deplVersionAndSecurityIpWhitelistFieldName:   "ip-whitelist",
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:        "oneshard",
				deplConfigurationNodeSizeIdFieldName:   "a8",
				deplConfigurationNodeCountFieldName:    3,
				deplConfigurationNodeDiskSizeFieldName: 32,
			},
		},
	}
	s := resourceDeployment().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedDepl := expandDeploymentResource(resourceData, "123456789")
	assert.Equal(t, depl, expandedDepl)
}

func TestExpandDeploymentOverrideProjectID(t *testing.T) {
	depl := &data.Deployment{
		Name:        "test-name",
		Description: "test-desc",
		ProjectId:   "overrideid",
		RegionId:    "gcp-europe-west4",
		Version:     "3.6.0",
		Certificates: &data.Deployment_CertificateSpec{
			CaCertificateId: "certificate-id",
		},
		IpwhitelistId: "ip-whitelist",
		Model: &data.Deployment_ModelSpec{
			Model:        "oneshard",
			NodeSizeId:   "a8",
			NodeCount:    3,
			NodeDiskSize: 32,
		},
	}
	raw := map[string]interface{}{
		deplProjectFieldName:     "overrideid",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionAndSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplVersionAndSecurityDbVersionFieldName:     "3.6.0",
				deplVersionAndSecurityCaCertificateFieldName: "certificate-id",
				deplVersionAndSecurityIpWhitelistFieldName:   "ip-whitelist",
			},
		},
		deplConfigurationFieldName: []interface{}{
			map[string]interface{}{
				deplConfigurationModelFieldName:        "oneshard",
				deplConfigurationNodeSizeIdFieldName:   "a8",
				deplConfigurationNodeCountFieldName:    3,
				deplConfigurationNodeDiskSizeFieldName: 32,
			},
		},
	}
	s := resourceDeployment().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedDepl := expandDeploymentResource(resourceData, "thisshouldbeoverriden")
	assert.Equal(t, depl, expandedDepl)
}
