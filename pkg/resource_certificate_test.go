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
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

var (
	testAccProviders   map[string]terraform.ResourceProvider
	testAccProvider    *schema.Provider
	testOrganizationId string
	testProject        *rm.Project
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"oasis": testAccProvider,
	}
	testOrganizationId = os.Getenv("OASIS_TEST_ORGANIZATION_ID")
	// Initialize Client with connection settings
	testAccProvider.Configure(terraform.NewResourceConfigRaw(nil))
}

func TestResourceCertificate_Basic(t *testing.T) {
	t.Parallel()
	res := "test-cert-" + acctest.RandString(10)
	name := "terraform-cert-" + acctest.RandString(10)
	id, err := getOrCreateProject()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestProject()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCertificate,
		Steps: []resource.TestStep{
			{
				Config: testBasicConfig(res, name, id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("oasis_certificate."+res, "description"),
					resource.TestCheckResourceAttr("oasis_certificate."+res, "name", name),
				),
			},
		},
	})
}

func getOrCreateProject() (string, error) {
	if testProject != nil {
		return testProject.GetId(), nil
	}

	client := testAccProvider.Meta().(*Client)
	err := client.Connect()
	if err != nil {
		return "", err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)

	proj, err := rmc.CreateProject(client.ctxWithToken, &rm.Project{
		Name:           "terraform-test-project",
		Description:    "This is a project used by terraform acceptance tests. PLEASE DO NOT DELETE!",
		OrganizationId: testOrganizationId,
	})
	if err != nil {
		return "", err
	}
	testProject = proj
	return testProject.GetId(), nil
}

func deleteTestProject() error {
	client := testAccProvider.Meta().(*Client)
	err := client.Connect()
	if err != nil {
		return err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	_, err = rmc.DeleteProject(client.ctxWithToken, &common.IDOptions{Id: testProject.GetId()})
	testProject = nil
	return err
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

func TestExpandingCertificateResource(t *testing.T) {
	raw := map[string]interface{}{
		"name":                       "test-name",
		"description":                "test-description",
		"project":                    "123456789",
		"use_well_known_certificate": true,
		"lifetime":                   3600,
	}
	s := resourceCertificate().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	cert := expandToCertificate(data)
	assert.Equal(t, raw["name"], cert.GetName())
	assert.Equal(t, raw["description"], cert.GetDescription())
	assert.Equal(t, raw["project"], cert.GetProjectId())
	assert.Equal(t, raw["use_well_known_certificate"], cert.GetUseWellKnownCertificate())
	assert.Equal(t, raw["lifetime"], int(cert.GetLifetime().GetSeconds()))
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

func testBasicConfig(resource, name, project string) string {
	return fmt.Sprintf(`resource "oasis_certificate" "%s" {
  name = "%s"
  description = "Terraform Updated Generated Certificate"
  project      = "%s"
  use_well_known_certificate = false
}`, resource, name, project)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
	if v := os.Getenv("OASIS_TEST_ORGANIZATION_ID"); v == "" {
		t.Fatal("the test needs an organization id to use for testing")
	}
}
