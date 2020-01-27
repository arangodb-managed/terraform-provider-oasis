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

func TestOasisProjectDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccDataSourcePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBasicOasisProjectDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.oasis_project.test", id),
					resource.TestCheckResourceAttrSet("data.oasis_project.test", name),
					resource.TestCheckResourceAttrSet("data.oasis_project.test", createdAt),
					resource.TestCheckResourceAttrSet("data.oasis_project.test", url),
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

func TestFlattenProjectDataSource(t *testing.T) {
	createdAtTimeStamp, _ := types.TimestampProto(time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC))
	proj := rm.Project{
		Id:             "test-id",
		Url:            "https://test.url",
		Name:           "test-name",
		Description:    "test-description",
		OrganizationId: "org-id",
		CreatedAt:      createdAtTimeStamp,
	}
	expected := map[string]interface{}{
		id:          "test-id",
		name:        "test-name",
		description: "test-description",
		url:         "https://test.url",
		createdAt:   "1980-01-01T01:01:01Z",
	}
	got := flattenProjectObject(&proj)
	assert.Equal(t, expected, got)
}

func testBasicOasisProjectDataSourceConfig() string {
	return fmt.Sprintf(`data "oasis_project" "test" {
	id = "168594080"
}`)
}
