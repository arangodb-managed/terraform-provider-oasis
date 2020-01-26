package pkg

import (
	"fmt"
	"os"
	"testing"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"arangodb": testAccProvider,
	}
}

func xTestResourceCertificate_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCertificate,
		Steps: []resource.TestStep{
			{
				Config: testBasicConfig(),
			},
			{
				ResourceName:      "test_oasis_certificate.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestFlattenCertificate(t *testing.T) {

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

		_, err := cryptoc.GetCACertificate(client.ctxWithToken, &common.IDOptions{Id: d.Id()})
		if err == nil {
			return fmt.Errorf("certificate still present")
		}
	}

	return nil
}

func testBasicConfig() string {
	return `resource "oasis_certificate" "my_oasis_cert" {
  name = "terraform-cert"
  description = "Terraform Updated Generated Certificate"
  project      = "168594080"
  use_well_known_certificate = false
}`
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}
