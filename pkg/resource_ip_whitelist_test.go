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
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	common "github.com/arangodb-managed/apis/common/v1"
	security "github.com/arangodb-managed/apis/security/v1"
)

func TestResourceIPWhitelist(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	res := "terraform-ipwhitelist-" + acctest.RandString(10)
	name := "ipwhitelist-" + acctest.RandString(10)
	orgID, err := FetchOrganizationID(testAccProvider)
	assert.NoError(t, err)
	pid, err := FetchProjectID(orgID, testAccProvider)
	assert.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyIPWhitelist,
		Steps: []resource.TestStep{
			{
				Config: testBasicConfig(res, name, pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_ipwhitelist."+res, ipCIDRRangeFieldName+".#", "3"),
					resource.TestCheckResourceAttr("oasis_ipwhitelist."+res, ipCIDRRangeFieldName+".0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr("oasis_ipwhitelist."+res, ipCIDRRangeFieldName+".1", "88.11.0.0/16"),
					resource.TestCheckResourceAttr("oasis_ipwhitelist."+res, ipCIDRRangeFieldName+".2", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("oasis_ipwhitelist."+res, ipNameFieldName, name),
				),
			},
		},
	})
}

func TestFlattenIPWhitelistResource(t *testing.T) {
	expected := map[string]interface{}{
		ipNameFieldName:        "test-name",
		ipDescriptionFieldName: "test-description",
		ipCreatedAtFieldName:   "1980-03-03T01:01:01Z",
		ipProjectFieldName:     "123456789",
		ipCIDRRangeFieldName:   []string{"1.2.3.4/32", "88.11.0.0/16", "0.0.0.0/0"},
		ipIsDeletedFieldName:   false,
	}

	created, _ := types.TimestampProto(time.Date(1980, 03, 03, 1, 1, 1, 0, time.UTC))
	cert := security.IPWhitelist{
		Name:        "test-name",
		Description: "test-description",
		ProjectId:   "123456789",
		CidrRanges:  []string{"1.2.3.4/32", "88.11.0.0/16", "0.0.0.0/0"},
		CreatedAt:   created,
		IsDeleted:   false,
	}
	got := flattenIPWhitelistResource(&cert)
	assert.Equal(t, expected, got)
}

func TestExpandingIPWhitelistResource(t *testing.T) {
	raw := map[string]interface{}{
		ipNameFieldName:        "test-name",
		ipDescriptionFieldName: "test-description",
		ipProjectFieldName:     "123456789",
		ipCIDRRangeFieldName:   []interface{}{"1.2.3.4/32", "88.11.0.0/16", "0.0.0.0/0"},
		ipIsDeletedFieldName:   false,
	}
	cidrRange, err := expandStringList(raw[ipCIDRRangeFieldName].([]interface{}))
	assert.NoError(t, err)
	s := resourceIPWhitelist().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	whitelist, err := expandToIPWhitelist(data, "123456789")
	assert.NoError(t, err)
	assert.Equal(t, raw[ipNameFieldName], whitelist.GetName())
	assert.Equal(t, raw[ipDescriptionFieldName], whitelist.GetDescription())
	assert.Equal(t, raw[ipIsDeletedFieldName], whitelist.GetIsDeleted())
	assert.Equal(t, raw[ipProjectFieldName], whitelist.GetProjectId())
	assert.Equal(t, cidrRange, whitelist.GetCidrRanges())
}

func TestExpandingIPWhitelistResourceNameNotDefinedError(t *testing.T) {
	raw := map[string]interface{}{
		ipDescriptionFieldName: "test-description",
		ipProjectFieldName:     "123456789",
		ipCIDRRangeFieldName:   []interface{}{"1.2.3.4/32", "88.11.0.0/16", "0.0.0.0/0"},
		ipIsDeletedFieldName:   false,
	}
	s := resourceIPWhitelist().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	_, err := expandToIPWhitelist(data, "123456789")
	assert.EqualError(t, err, "failed to parse field "+ipNameFieldName)
}

func testAccCheckDestroyIPWhitelist(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	securityc := security.NewSecurityServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_ipwhitelist" {
			continue
		}

		if _, err := securityc.DeleteIPWhitelist(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); err == nil {
			return fmt.Errorf("IPWhitelist still present")
		}
	}

	return nil
}

func testBasicConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_ipwhitelist" "%s" {
  name = "%s"
  description = "Terraform Generated IPWhitelist"
  project      = "%s"
  cidr_ranges = ["1.2.3.4/32", "88.11.0.0/16", "0.0.0.0/0"]
}`, resource, name, project)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}
