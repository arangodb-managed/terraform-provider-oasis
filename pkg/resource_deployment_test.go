//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
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
//

package pkg

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	data "github.com/arangodb-managed/apis/data/v1"
)

func TestFlattenDeploymentResource(t *testing.T) {
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
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.6.0",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName: "certificate-id",
				deplSecurityIpWhitelistFieldName:   "ip-whitelist",
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
	raw := map[string]interface{}{
		deplProjectFieldName:     "123456789",
		deplNameFieldName:        "test-name",
		deplDescriptionFieldName: "test-desc",
		deplLocationFieldName: []interface{}{
			map[string]interface{}{
				deplLocationRegionFieldName: "gcp-europe-west4",
			},
		},
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.6.0",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName: "certificate-id",
				deplSecurityIpWhitelistFieldName:   "ip-whitelist",
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
	expandedDepl, err := expandDeploymentResource(resourceData, "123456789")
	assert.NoError(t, err)
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
		deplVersionFieldName: []interface{}{
			map[string]interface{}{
				deplVersionDbVersionFieldName: "3.6.0",
			},
		},
		deplSecurityFieldName: []interface{}{
			map[string]interface{}{
				deplSecurityCaCertificateFieldName: "certificate-id",
				deplSecurityIpWhitelistFieldName:   "ip-whitelist",
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
	expandedDepl, err := expandDeploymentResource(resourceData, "thisshouldbeoverriden")
	assert.NoError(t, err)
	assert.Equal(t, depl, expandedDepl)
}
