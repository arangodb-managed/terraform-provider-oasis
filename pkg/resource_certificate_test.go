package pkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"oasis": testAccProvider,
	}
}

func TestResourceCertificate_Basic(t *testing.T) {
	t.Parallel()
	res := "test-cert-" + acctest.RandString(10)
	name := "terraform-cert-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCertificate,
		Steps: []resource.TestStep{
			{
				Config: testBasicConfig(res, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("oasis_certificate."+res, "description"),
					resource.TestCheckResourceAttr("oasis_certificate."+res, "name", name),
				),
			},
		},
	})
}

func TestFlattenCertificateResource(t *testing.T) {
	expected := map[string]interface{}{
		"name":                       "test-name",
		"description":                "test-description",
		"project":                    "123456789",
		"use_well_known_certificate": true,
		"lifetime":                   3600,
		"is_default":                 false,
		"expires_at":                 "1980-03-10T01:01:01Z",
		"created_at":                 "1980-03-03T01:01:01Z",
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

func testAccCheckDestroyCertificate(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	err := client.Connect()
	if err != nil {
		return err
	}
	cryptoc := crypto.NewCryptoServiceClient(client.conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oasis_certificate" {
			continue
		}

		_, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: rs.Primary.ID})
		if err == nil {
			return fmt.Errorf("certificate still present")
		}
	}

	return nil
}

func testBasicConfig(resource, name string) string {
	return fmt.Sprintf(`resource "oasis_certificate" "%s" {
  name = "%s"
  description = "Terraform Updated Generated Certificate"
  project      = "168594080"
  use_well_known_certificate = false
}`, resource, name)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}
