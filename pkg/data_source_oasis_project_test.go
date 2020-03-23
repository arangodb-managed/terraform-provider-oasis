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

func TestOasisProjectDataSource_Basic(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	if _, ok := os.LookupEnv("OASIS_TEST_ORGANIZATION_ID"); !ok {
		t.Skip("This test requires an organization id to be set.")
	}
	pid, err := fetchProjectID()
	if err != nil {
		t.Fatal(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testProjectAccDataSourcePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBasicOasisProjectDataSourceConfig(pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.oasis_project.test", projIdFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_project.test", projNameFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_project.test", projCreatedAtFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_project.test", projUrlFieldName),
				),
			},
		},
	})
}

// fetchProjectID finds and retrieves the first project ID it finds in the given Organization.
func fetchProjectID() (string, error) {
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
	orgID := os.Getenv("OASIS_TEST_ORGANIZATION_ID")
	if proj, err := rmc.ListProjects(client.ctxWithToken, &common.ListOptions{ContextId: orgID}); err != nil {
		client.log.Error().Err(err).Msg("Failed to list projects")
		return "", err
	} else if len(proj.Items) < 1 {
		client.log.Error().Err(err).Msg("No projects found in organization")
		return "", fmt.Errorf("no projects found in organization %s", orgID)
	} else {
		return proj.Items[0].GetId(), nil
	}
}

func testProjectAccDataSourcePreCheck(t *testing.T) {
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
		projIdFieldName:          "test-id",
		projNameFieldName:        "test-name",
		projDescriptionFieldName: "test-description",
		projUrlFieldName:         "https://test.url",
		projCreatedAtFieldName:   "1980-01-01T01:01:01Z",
	}
	got := flattenProjectObject(&proj)
	assert.Equal(t, expected, got)
}

func testBasicOasisProjectDataSourceConfig(pid string) string {
	return fmt.Sprintf(`data "oasis_project" "test" {
	id = "%s"
}`, pid)
}
