//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

func TestOasisOrganizationDataSource_Basic(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	organizationID, err := FetchOrganizationID(testAccProvider)
	if err != nil {
		t.Fatal(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccDataSourcePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBasicOasisOrganizationDataSourceConfig(organizationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", idFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", nameFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", createdAtFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_organization.test", urlFieldName),
				),
			},
		},
	})
}

func testAccDataSourcePreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}

func TestFlattenOrganizationDataSource(t *testing.T) {
	createdAtTimeStamp, _ := types.TimestampProto(time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC))
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
		idFieldName:          "test-id",
		nameFieldName:        "test-name",
		descriptionFieldName: "test-description",
		urlFieldName:         "https://test.url",
		createdAtFieldName:   "1980-01-01T01:01:01Z",
		tierFieldName:        flattenedTier,
	}
	got := flattenOrganizationObject(&org)
	assert.Equal(t, expected[idFieldName], got[idFieldName])
	assert.Equal(t, expected[nameFieldName], got[nameFieldName])
	assert.Equal(t, expected[descriptionFieldName], got[descriptionFieldName])
	assert.Equal(t, expected[urlFieldName], got[urlFieldName])
	assert.Equal(t, expected[createdAtFieldName], got[createdAtFieldName])
	assert.True(t, flattenedTier.Equal(got[tierFieldName]))

}

func testBasicOasisOrganizationDataSourceConfig(id string) string {
	return fmt.Sprintf(`data "oasis_organization" "test" {
	id = "%s"
}`, id)
}
