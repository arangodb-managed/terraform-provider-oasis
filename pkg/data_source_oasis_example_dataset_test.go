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
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"

	example "github.com/arangodb-managed/apis/example/v1"
)

func TestFlattenExampleDataset(t *testing.T) {
	created, _ := types.TimestampProto(time.Date(1980, 03, 03, 1, 1, 1, 0, time.UTC))
	items := &example.ExampleDatasetList{
		Items: []*example.ExampleDataset{
			{
				Id:          "test-id",
				Url:         "test-url",
				CreatedAt:   created,
				Guide:       "test-guide",
				Name:        "test-name",
				Description: "test-description",
			},
		},
	}
	expected := map[string]interface{}{
		exampleOrganizationIDFieldName: "",
		exampleExampleDatasetsFieldName: []interface{}{
			map[string]interface{}{
				exampleExampleDatasetsIDFieldName:          "test-id",
				exampleExampleDatasetsNameFieldName:        "test-name",
				exampleExampleDatasetsDescriptionFieldName: "test-description",
				exampleExampleDatasetsGuideFieldName:       "test-guide",
				exampleExampleDatasetsCreatedAtFieldName:   "1980-03-03T01:01:01Z",
			},
		},
	}
	flattened := flattenExampleDatasets("", items.Items)
	assert.Equal(t, expected, flattened)
}
