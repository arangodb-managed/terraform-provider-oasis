//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestOasisTermsAndConditionsDataSource_Basic(t *testing.T) {
	if _, ok := os.LookupEnv("TF_ACC"); !ok {
		t.Skip()
	}
	if _, ok := os.LookupEnv("TF_OASIS_TERMS_AND_CONDITIONS"); !ok {
		t.Skip()
	}
	tAndCID := os.Getenv("TF_OASIS_TERMS_AND_CONDITIONS")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testTAndCAccDataSourcePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBasicOasisTandCDataSourceConfig(tAndCID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.oasis_terms_and_conditions.test_t_and_c", tcIDFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_terms_and_conditions.test_t_and_c", tcContentFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_terms_and_conditions.test_t_and_c", tcCreatedAtFieldName),
					resource.TestCheckResourceAttrSet("data.oasis_terms_and_conditions.test_t_and_c", tcOrganizationFieldName),
				),
			},
		},
	})
}

func TestFlattenTAndCDataSource(t *testing.T) {
	createdAtTimeStamp, _ := types.TimestampProto(time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC))
	term := rm.TermsAndConditions{
		Id:        "test-id",
		Content:   "Content",
		CreatedAt: createdAtTimeStamp,
	}
	expected := map[string]interface{}{
		tcIDFieldName:        "test-id",
		tcCreatedAtFieldName: "1980-01-01T01:01:01Z",
		tcContentFieldName:   "Content",
	}
	got := flattenTCObject(&term)
	assert.Equal(t, expected, got)
}

func testTAndCAccDataSourcePreCheck(t *testing.T) {
	if v := os.Getenv("OASIS_API_KEY_ID"); v == "" {
		t.Fatal("the test needs a test account key to run")
	}
}

func testBasicOasisTandCDataSourceConfig(pid string) string {
	return fmt.Sprintf(`data "oasis_terms_and_conditions" "test_t_and_c" {
	id = "%s"
}`, pid)
}
