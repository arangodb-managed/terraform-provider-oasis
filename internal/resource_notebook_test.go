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

package provider

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/types"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	nb "github.com/arangodb-managed/apis/notebook/v1"
)

// TestAccResourceNotebook verifies the Oasis Notebook resource is created along with the specified properties
func TestAccResourceNotebook(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	resourceName := "terraform-notebook-" + acctest.RandString(10)

	orgID, err := FetchOrganizationID()
	require.NoError(t, err)
	projectID, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyNotebook,
		Steps: []resource.TestStep{
			{
				Config:      testNotebookConfig("", resourceName),
				ExpectError: regexp.MustCompile("InvalidArgument desc = Project ID missing"),
			},
			{
				Config: testNotebookConfig(projectID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_deployment.my_oneshard_deployment", deplNameFieldName, "oasis_notebook_deployment"),
					resource.TestCheckResourceAttr("oasis_notebook."+resourceName, notebookNameFieldName, "Test-Notebook"),
				),
			},
		},
	})
}

// testNotebookConfig contains the Terraform resource definitions for testing usage
func testNotebookConfig(project, notebookResource string) string {
	return fmt.Sprintf(`
	resource "oasis_deployment" "my_oneshard_deployment" {
		terms_and_conditions_accepted = "true"
		project = "%s"
		name = "oasis_notebook_deployment"
		location {
			region = "gcp-europe-west4"
		}
		security {
			disable_foxx_authentication = false
		}
		disk_performance = "dp30"
		configuration {
			model = "oneshard"
			node_size_id = "c4-a8"
			node_disk_size = 20
			maximum_node_disk_size = 40
		}
		notification_settings {
			email_addresses = [
			"test@arangodb.com"
			]
		}
	}

	resource "oasis_notebook" "%s" {
		deployment_id = oasis_deployment.my_oneshard_deployment.id
		name          = "Test-Notebook"
		model {
			notebook_model_id = "basic"
			disk_size         = "10"
  		}
	}
`, project, notebookResource)
}

// TestFlattenOrganization tests the Oasis Organization flattening for Terraform schema compatibility.
func TestFlattenNotebook(t *testing.T) {
	created, _ := types.TimestampProto(time.Date(2022, 03, 03, 1, 1, 1, 0, time.UTC))

	notebook := &nb.Notebook{
		Id:           "test",
		DeploymentId: "axt9evhsotaxtfnk9qml",
		Name:         "Test-Notebook",
		Description:  "Jupyter Notebook description",

		Model: &nb.ModelSpec{
			NotebookModelId: "taxtxt9evhsofnk9qmla",
			DiskSize:        20,
		},
		CreatedAt: created,
	}

	expected := map[string]interface{}{
		notebookDeploymentIdFieldName: "axt9evhsotaxtfnk9qml",
		notebookNameFieldName:         "Test-Notebook",
		notebookDescriptionFieldName:  "Jupyter Notebook description",
		notebookIsPausedFieldName:     false,
		notebookIsDeletedFieldName:    false,
		notebookCreatedAtFieldName:    "2022-03-03T01:01:01Z",
		notebookModelFieldName: []interface{}{
			map[string]interface{}{
				notebookModelIdFieldName:       "taxtxt9evhsofnk9qmla",
				notebookModelDiskSizeFieldName: int32(20),
			},
		},
	}
	flattened := flattenNotebookResource(notebook)
	assert.Equal(t, expected, flattened)
}

// TestExpandNotebook tests the Oasis Notebook expansion for Terraform schema compatibility.
func TestExpandNotebook(t *testing.T) {
	raw := map[string]interface{}{
		notebookDeploymentIdFieldName: "axt9evhsotaxtfnk9qml",
		notebookNameFieldName:         "Test-Notebook",
		notebookDescriptionFieldName:  "Test description for Jupyter Notebook",
		notebookModelFieldName: []interface{}{
			map[string]interface{}{
				notebookModelIdFieldName:       "taxtxt9evhsofnk9qmla",
				notebookModelDiskSizeFieldName: 20,
			},
		},
	}
	expected := &nb.Notebook{
		DeploymentId: "axt9evhsotaxtfnk9qml",
		Name:         "Test-Notebook",
		Description:  "Test description for Jupyter Notebook",
		Model: &nb.ModelSpec{
			NotebookModelId: "taxtxt9evhsofnk9qmla",
			DiskSize:        20,
		},
	}

	s := resourceNotebook().Schema
	resourceData := schema.TestResourceDataRaw(t, s, raw)
	expandedNotebook, err := expandNotebookResource(resourceData)
	assert.NoError(t, err)

	assert.Equal(t, expected, expandedNotebook)
}

// testAccCheckDestroyNotebook verifies the Terraform oasis_notebook resource cleanup.
func testAccCheckDestroyNotebook(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	nbc := nb.NewNotebookServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_notebook" {
			continue
		}

		if _, err := nbc.GetNotebook(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); !common.IsNotFound(err) {
			return fmt.Errorf("notebook still present")
		}
	}

	return nil
}
