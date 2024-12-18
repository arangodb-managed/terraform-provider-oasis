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
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

func TestOasisProjectDataSource_Basic(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	orgID, err := FetchOrganizationID()
	require.NoError(t, err)
	pid, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	if err != nil {
		t.Fatal(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testProjectAccDataSourcePreCheck(t) },
		ProviderFactories: testProviderFactories,
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

func testProjectAccDataSourcePreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}

func TestFlattenProjectDataSource(t *testing.T) {
	createdAtTimeStamp := timestamppb.New(time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC))
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
