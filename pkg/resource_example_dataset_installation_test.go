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
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	example "github.com/arangodb-managed/apis/example/v1"
)

func TestFlattenExampleDatasetInstallation(t *testing.T) {
	created, _ := types.TimestampProto(time.Date(1980, 03, 03, 1, 1, 1, 0, time.UTC))
	testInt := &example.ExampleDatasetInstallation{
		Id:               "test-id",
		Url:              "test-url",
		DeploymentId:     "test-depl-id",
		ExampledatasetId: "test-example-id",
		CreatedAt:        created,
		Status: &example.ExampleDatasetInstallation_Status{
			DatabaseName: "test-db",
			State:        example.ExampleInstallationStatusReady,
			IsFailed:     false,
			IsAvailable:  true,
		},
	}
	expected := map[string]interface{}{
		datasetDeploymentIdFieldName:     "test-depl-id",
		datasetExampleDatasetIdFieldName: "test-example-id",
		datasetCreatedAtFieldName:        "1980-03-03T01:01:01Z",
		datasetStatusFieldName: []interface{}{
			map[string]interface{}{
				datasetStatusStateFieldName:        example.ExampleInstallationStatusReady,
				datasetStatusIsFailedFieldName:     false,
				datasetStatusIsAvailableFieldName:  true,
				datasetStatusDatabaseNameFieldName: "test-db",
			},
		},
	}
	flattened := flattenExampleDatasetInstallation(testInt)
	assert.Equal(t, expected, flattened)
}

func TestExpandExampleDatasetInstallation(t *testing.T) {
	raw := map[string]interface{}{
		datasetDeploymentIdFieldName:     "test-depl-id",
		datasetExampleDatasetIdFieldName: "test-example-id",
	}
	s := resourceExampleDatasetInstallation().Schema
	data := schema.TestResourceDataRaw(t, s, raw)
	expanded := expandExampleDatasetInstallation(data)
	assert.Equal(t, raw[datasetDeploymentIdFieldName], expanded.DeploymentId)
	assert.Equal(t, raw[datasetExampleDatasetIdFieldName], expanded.ExampledatasetId)
}
