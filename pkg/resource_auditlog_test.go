//
// DISCLAIMER
//
// Copyright 2022 ArangoDB GmbH, Cologne, Germany
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

package pkg

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	audit "github.com/arangodb-managed/apis/audit/v1"
)

// TestFlattenAuditLog tests the Oasis Audit Log flattening for Terraform schema compatibility.
func TestFlattenAuditLog(t *testing.T) {
	auditLog := &audit.AuditLog{
		Name:           "test-auditlog",
		Description:    "test-description",
		OrganizationId: "9047335679",
	}

	expected := map[string]interface{}{
		auditLogNameFieldName:         "test-auditlog",
		auditLogDescriptionFieldName:  "test-description",
		auditLogOrganizationFieldName: "9047335679",
	}

	flattened := flattenAuditLogResource(auditLog)
	assert.Equal(t, expected, flattened)
}

// TestExpandAuditLog tests the Oasis AuditLog expansion for Terraform schema compatibility.
func TestExpandAuditLog(t *testing.T) {
	raw := map[string]interface{}{
		auditLogNameFieldName:         "test-auditlog",
		auditLogDescriptionFieldName:  "test-description",
		auditLogOrganizationFieldName: "9047335679",
	}

	expected := &audit.AuditLog{
		Name:           "test-auditlog",
		Description:    "test-description",
		OrganizationId: "9047335679",
		Destinations:   []*audit.AuditLog_Destination{{Type: audit.DestinationCloud}},
	}

	s := resourceAuditLog().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedAuditLog, err := expandAuditLogResource(resourceData)
	assert.NoError(t, err)

	assert.Equal(t, expected, expandedAuditLog)
}
