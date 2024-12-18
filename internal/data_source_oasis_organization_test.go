//
// DISCLAIMER
//
// Copyright 2020-2024 ArangoDB GmbH, Cologne, Germany
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

package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

func TestOasisOrganizationDataSource_Basic(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	organizationID, err := FetchOrganizationID()
	if err != nil {
		t.Fatal(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testOrgAccDataSourcePreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBasicOasisOrganizationDataSourceConfig(organizationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", orgIdFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", orgNameFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", orgCreatedAtFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", orgUrlFieldName),
				),
			},
		},
	})
}

func testOrgAccDataSourcePreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}

func TestFlattenOrganizationDataSource(t *testing.T) {
	createdAtTimeStamp := timestamppb.New(time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC))
	org := rm.Organization{
		Id:          "test-id",
		Url:         "https://test.url",
		Name:        "test-name",
		Description: "test-description",
		CreatedAt:   createdAtTimeStamp,
		Tier: &rm.Tier{
			Id:                         "free",
			Name:                       "Free to try",
			HasSupportPlans:            true,
			HasBackupUploads:           true,
			RequiresTermsAndConditions: true,
		},
	}
	flattenedTier := flattenTierObject(org.Tier)
	expected := map[string]interface{}{
		orgIdFieldName:          "test-id",
		orgNameFieldName:        "test-name",
		orgDescriptionFieldName: "test-description",
		orgUrlFieldName:         "https://test.url",
		orgCreatedAtFieldName:   "1980-01-01T01:01:01Z",
		tierFieldName:           flattenedTier,
	}
	got := flattenOrganizationObject(&org)
	assert.Equal(t, expected[orgIdFieldName], got[orgIdFieldName])
	assert.Equal(t, expected[orgNameFieldName], got[orgNameFieldName])
	assert.Equal(t, expected[orgDescriptionFieldName], got[orgDescriptionFieldName])
	assert.Equal(t, expected[orgUrlFieldName], got[orgUrlFieldName])
	assert.Equal(t, expected[orgCreatedAtFieldName], got[orgCreatedAtFieldName])
	assert.True(t, flattenedTier.Equal(got[tierFieldName]))

}

func testBasicOasisOrganizationDataSourceConfig(id string) string {
	return fmt.Sprintf(`data "oasis_organization" "test" {
	id = "%s"
}`, id)
}
