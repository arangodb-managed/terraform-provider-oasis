//
// DISCLAIMER
//
// Copyright 2020-2022 ArangoDB GmbH, Cologne, Germany
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
//

package pkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

func TestResourceCreateProject(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	res := "terraform-project-" + acctest.RandString(10)
	name := "terraform-project-name" + acctest.RandString(10)
	orgID, err := FetchOrganizationID()
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: testBasicProjectConfig(res, name, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_project."+res, projectNameFieldName, name),
				),
			},
		},
	})
}

func TestFlattenProjectResource(t *testing.T) {
	expected := map[string]interface{}{
		projectNameFieldName:         "test-name",
		projectDescriptionFieldName:  "test-description",
		projectCreatedAtFieldName:    "1980-03-03T01:01:01Z",
		projectOrganizationFieldName: "_support",
		projectIsDeletedFieldName:    true,
	}

	created, _ := types.TimestampProto(time.Date(1980, 03, 03, 1, 1, 1, 0, time.UTC))
	proj := rm.Project{
		Name:           "test-name",
		Description:    "test-description",
		OrganizationId: "_support",
		CreatedAt:      created,
		IsDeleted:      true,
	}
	got := flattenProjectResource(&proj)
	assert.Equal(t, expected, got)
}

func TestExpandingProjectResource(t *testing.T) {
	raw := map[string]interface{}{
		projectNameFieldName:         "test-name",
		projectDescriptionFieldName:  "test-description",
		projectOrganizationFieldName: "_support",
	}
	s := resourceProject().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	project, err := expandToProject(data, "123456789")
	assert.NoError(t, err)
	assert.Equal(t, raw[projectNameFieldName], project.GetName())
	assert.Equal(t, raw[projectDescriptionFieldName], project.GetDescription())
	assert.Equal(t, raw[projectOrganizationFieldName], project.GetOrganizationId())
}

func testAccCheckDestroyProject(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_project" {
			continue
		}

		if _, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); !common.IsNotFound(err) {
			return fmt.Errorf("project still present")
		}
	}

	return nil
}

func testBasicProjectConfig(res string, name string, id string) string {
	return fmt.Sprintf(`resource "oasis_project" "%s" {
  name = "%s"
  description = "Terraform Generated Project"
  organization = "%s"
}`, res, name, id)
}
