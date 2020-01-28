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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

func TestOasisOrganizationDataSource_Basic(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	organizationID, err := fetchOrganizationID()
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

// fetchOrganizationID finds and retrieves the first Organization ID it finds in the given Organization.
func fetchOrganizationID() (string, error) {
	// Initialize Client with connection settings
	if err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil)); err != nil {
		return "", err
	}
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return "", err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	if organizations, err := rmc.ListOrganizations(client.ctxWithToken, &common.ListOptions{}); err != nil {
		client.log.Error().Err(err).Msg("Failed to list Organizations")
		return "", err
	} else if len(organizations.GetItems()) < 1 {
		client.log.Error().Err(err).Msg("No Organizations found")
		return "", fmt.Errorf("no organizations found")
	} else {
		return organizations.GetItems()[0].GetId(), nil
	}
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
	expected := map[string]interface{}{
		idFieldName:          "test-id",
		nameFieldName:        "test-name",
		descriptionFieldName: "test-description",
		urlFieldName:         "https://test.url",
		createdAtFieldName:   "1980-01-01T01:01:01Z",
		tierFieldName: map[string]interface{}{
			tierIdFieldName:                         "free",
			tierNameFieldName:                       "Free to try",
			tierHasSupportPlansFieldName:            "true",
			tierHasBackupUploadsFieldName:           "true",
			tierRequiresTermsAndConditionsFieldName: "true",
		},
	}
	got := flattenOrganizationObject(&org)
	assert.Equal(t, expected, got)
}

func testBasicOasisOrganizationDataSourceConfig(id string) string {
	return fmt.Sprintf(`data "oasis_organization" "test" {
	id = "%s"
}`, id)
}
