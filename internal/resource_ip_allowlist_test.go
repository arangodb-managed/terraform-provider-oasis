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
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	common "github.com/arangodb-managed/apis/common/v1"
	security "github.com/arangodb-managed/apis/security/v1"
)

func TestResourceIPAllowlist(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()

	res := "terraform-ipallowlist-" + acctest.RandString(10)
	name := "ipallowlist-" + acctest.RandString(10)
	orgID, err := FetchOrganizationID()
	assert.NoError(t, err)
	pid, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	assert.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyIPAllowlist,
		Steps: []resource.TestStep{
			{
				Config:      testBasicConfig(res, "", pid),
				ExpectError: regexp.MustCompile("failed to parse field name"),
			},
			{
				Config: testLockedIpAllowListConfig(res, name, pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_ipallowlist."+res, ipLockedFieldName, "true"),
				),
			},
			{
				Config: testBasicConfig(res, name, pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_ipallowlist."+res, ipCIDRRangeFieldName+".#", "2"),
					resource.TestCheckResourceAttr("oasis_ipallowlist."+res, ipCIDRRangeFieldName+".0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr("oasis_ipallowlist."+res, ipCIDRRangeFieldName+".1", "88.11.0.0/16"),
					resource.TestCheckResourceAttr("oasis_ipallowlist."+res, ipNameFieldName, name),
				),
			},
		},
	})
}

func TestFlattenIPAllowlistResource(t *testing.T) {
	expected := map[string]interface{}{
		ipNameFieldName:                    "test-name",
		ipDescriptionFieldName:             "test-description",
		ipCreatedAtFieldName:               "1980-03-03T01:01:01Z",
		ipProjectFieldName:                 "123456789",
		ipCIDRRangeFieldName:               []string{"1.2.3.4/32", "88.11.0.0/16"},
		ipRemoteInspectionAllowedFieldName: false,
		ipIsDeletedFieldName:               false,
		ipLockedFieldName:                  true,
	}

	created := timestamppb.New(time.Date(1980, 03, 03, 1, 1, 1, 0, time.UTC))
	cert := security.IPAllowlist{
		Name:        "test-name",
		Description: "test-description",
		ProjectId:   "123456789",
		CidrRanges:  []string{"1.2.3.4/32", "88.11.0.0/16"},
		CreatedAt:   created,
		IsDeleted:   false,
		Locked:      true,
	}
	got := flattenIPAllowlistResource(&cert)
	assert.Equal(t, expected, got)

	cert.RemoteInspectionAllowed = true
	expected[ipRemoteInspectionAllowedFieldName] = true
	got = flattenIPAllowlistResource(&cert)
	assert.Equal(t, expected, got)

	cert.RemoteInspectionAllowed = false
	expected[ipRemoteInspectionAllowedFieldName] = false
	got = flattenIPAllowlistResource(&cert)
	assert.Equal(t, expected, got)
}

func TestExpandingIPAllowlistResource(t *testing.T) {
	raw := map[string]interface{}{
		ipNameFieldName:                    "test-name",
		ipDescriptionFieldName:             "test-description",
		ipProjectFieldName:                 "123456789",
		ipCIDRRangeFieldName:               []interface{}{"1.2.3.4/32", "88.11.0.0/16"},
		ipRemoteInspectionAllowedFieldName: true,
		ipIsDeletedFieldName:               false,
		ipLockedFieldName:                  true,
	}
	cidrRange, err := expandStringList(raw[ipCIDRRangeFieldName].([]interface{}))
	assert.NoError(t, err)
	s := resourceIPAllowlist().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	allowlist, err := expandToIPAllowlist(data, "123456789")
	assert.NoError(t, err)
	assert.Equal(t, raw[ipNameFieldName], allowlist.GetName())
	assert.Equal(t, raw[ipDescriptionFieldName], allowlist.GetDescription())
	assert.Equal(t, raw[ipIsDeletedFieldName], allowlist.GetIsDeleted())
	assert.Equal(t, raw[ipProjectFieldName], allowlist.GetProjectId())
	assert.Equal(t, raw[ipRemoteInspectionAllowedFieldName], allowlist.GetRemoteInspectionAllowed())
	assert.Equal(t, raw[ipLockedFieldName], allowlist.GetLocked())
	assert.Equal(t, cidrRange, allowlist.GetCidrRanges())
}

func TestExpandingIPAllowlistResourceNameNotDefinedError(t *testing.T) {
	raw := map[string]interface{}{
		ipDescriptionFieldName: "test-description",
		ipProjectFieldName:     "123456789",
		ipCIDRRangeFieldName:   []interface{}{"1.2.3.4/32", "88.11.0.0/16"},
		ipIsDeletedFieldName:   false,
		ipLockedFieldName:      true,
	}
	s := resourceIPAllowlist().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	_, err := expandToIPAllowlist(data, "123456789")
	assert.EqualError(t, err, "failed to parse field "+ipNameFieldName)
}

func testAccCheckDestroyIPAllowlist(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	securityc := security.NewSecurityServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_ipallowlist" {
			continue
		}

		if _, err := securityc.GetIPAllowlist(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); !common.IsNotFound(err) {
			return fmt.Errorf("IPAllowlist still present")
		}
	}

	return nil
}

func testBasicConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_ipallowlist" "%s" {
  name         = "%s"
  description  = "Terraform Generated IPAllowlist"
  project      = "%s"
  cidr_ranges  = ["1.2.3.4/32", "88.11.0.0/16"]
}`, resource, name, project)
}

func testLockedIpAllowListConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_ipallowlist" "%s" {
  name         = "%s"
  description  = "Terraform Generated IPAllowlist"
  project      = "%s"
  cidr_ranges  = ["1.2.3.4/32", "88.11.0.0/16"]
  locked       = true
}`, resource, name, project)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}
