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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
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
		PreCheck:          func() { testTAndCAccDataSourcePreCheck(t) },
		ProviderFactories: testProviderFactories,
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
