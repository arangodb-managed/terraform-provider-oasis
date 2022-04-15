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

package pkg

import (
	"context"
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

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
)

func TestResourceCertificate(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	t.Parallel()
	res := "test-cert-" + acctest.RandString(10)
	certName := "terraform-cert-" + acctest.RandString(10)
	orgID, err := FetchOrganizationID()
	if err != nil {
		t.Fatal(err)
	}
	pid, err := FetchProjectID(context.Background(), orgID, testAccProvider)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCertificatePreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCheckDestroyCertificate,
		Steps: []resource.TestStep{
			{
				Config: testBasicCertificateConfig(res, certName, pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("oasis_certificate."+res, descriptionFieldName),
					resource.TestCheckResourceAttr("oasis_certificate."+res, nameFieldName, certName),
				),
			},
			{
				Config: testUseWellKnownConfig(res, certName, pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("oasis_certificate."+res, descriptionFieldName),
					resource.TestCheckResourceAttr("oasis_certificate."+res, nameFieldName, certName),
					resource.TestCheckResourceAttr("oasis_certificate."+res, useWellKnownCertificateFieldName, "true"),
				),
			},
			{
				Config: testOptionalFieldsConfig(res, certName, pid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oasis_certificate."+res, nameFieldName, certName),
				),
			},
		},
	})
}

func TestFlattenCertificateResource(t *testing.T) {
	expected := map[string]interface{}{
		nameFieldName:                    "test-name",
		descriptionFieldName:             "test-description",
		projectFieldName:                 "123456789",
		useWellKnownCertificateFieldName: true,
		lifetimeFieldName:                3600,
		isDefaultFieldName:               false,
		expiresAtFieldName:               "1980-03-10T01:01:01Z",
		createdAtFieldName:               "1980-03-03T01:01:01Z",
	}

	created, _ := types.TimestampProto(time.Date(1980, 03, 03, 1, 1, 1, 0, time.UTC))
	expires, _ := types.TimestampProto(time.Date(1980, 03, 10, 1, 1, 1, 0, time.UTC))
	cert := crypto.CACertificate{
		Name:                    "test-name",
		Description:             "test-description",
		ProjectId:               "123456789",
		Lifetime:                types.DurationProto(1 * time.Hour),
		CreatedAt:               created,
		ExpiresAt:               expires,
		IsDefault:               false,
		UseWellKnownCertificate: true,
	}
	got := flattenCertificateResource(&cert)
	assert.Equal(t, expected, got)
}

func TestExpandingCertificateResource(t *testing.T) {
	raw := map[string]interface{}{
		nameFieldName:                    "test-name",
		descriptionFieldName:             "test-description",
		projectFieldName:                 "123456789",
		useWellKnownCertificateFieldName: true,
		lifetimeFieldName:                3600,
	}
	s := resourceCertificate().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	cert := expandToCertificate(data)
	assert.Equal(t, raw[nameFieldName], cert.GetName())
	assert.Equal(t, raw[descriptionFieldName], cert.GetDescription())
	assert.Equal(t, raw[projectFieldName], cert.GetProjectId())
	assert.Equal(t, raw[useWellKnownCertificateFieldName], cert.GetUseWellKnownCertificate())
	assert.Equal(t, raw[lifetimeFieldName], int(cert.GetLifetime().GetSeconds()))
}

func testAccCheckDestroyCertificate(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}
	cryptoc := crypto.NewCryptoServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_certificate" {
			continue
		}

		if _, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID}); err == nil {
			return fmt.Errorf("certificate still present")
		}
	}

	return nil
}

func testBasicCertificateConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_certificate" "%s" {
  name = "%s"
  description = "Terraform Updated Generated Certificate"
  project      = "%s"
  use_well_known_certificate = false
}`, resource, name, project)
}

func testUseWellKnownConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_certificate" "%s" {
  name = "%s"
  description = "Terraform Updated Generated Certificate"
  project      = "%s"
  use_well_known_certificate = true
}`, resource, name, project)
}

func testOptionalFieldsConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_certificate" "%s" {
  name = "%s"
  project      = "%s"
}`, resource, name, project)
}

func testAccCertificatePreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}
