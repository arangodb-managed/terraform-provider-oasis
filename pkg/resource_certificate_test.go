package pkg

import (
	"fmt"
	"os"
	"testing"

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
	name := "terraform-cert-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCertificate,
		Steps: []resource.TestStep{
			{
				Config: testBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("oasis_certificate.my_oasis_cert", "description"),
					resource.TestCheckResourceAttr("oasis_certificate.my_oasis_cert", "name", name),
				),
			},
		},
	})
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

func testBasicConfig(name string) string {
	return fmt.Sprintf(`resource "oasis_certificate" "my_oasis_cert" {
  name = "%s"
  description = "Terraform Updated Generated Certificate"
  project      = "168594080"
  use_well_known_certificate = false
}`, name)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}
