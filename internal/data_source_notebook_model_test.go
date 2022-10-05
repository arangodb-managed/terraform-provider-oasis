//
// DISCLAIMER
//
// Copyright 2022 ArangoDB GmbH, Cologne, Germany
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
	"testing"

	"github.com/stretchr/testify/assert"

	nb "github.com/arangodb-managed/apis/notebook/v1"
)

func TestFlattenNotebookModels(t *testing.T) {
	items := &nb.NotebookModelList{
		Items: []*nb.NotebookModel{
			{
				Id:          "test-id",
				Name:        "test-name",
				Cpu:         2,
				Memory:      10,
				MaxDiskSize: 20,
				MinDiskSize: 10,
			},
		},
	}
	expected := map[string]interface{}{
		notebookModelDataSourceDeploymentIdFieldName: "test-depl-id",
		notebookModelDataSourceItemsFieldName: []interface{}{
			map[string]interface{}{
				notebookModelDataSourceIdFieldName:          "test-id",
				notebookModelDataSourceNameFieldName:        "test-name",
				notebookModelDataSourceCPUFieldName:         float32(2),
				notebookModelDataSourceMemoryFieldName:      int32(10),
				notebookModelDataSourceMaxDiskSizeFieldName: int32(20),
				notebookModelDataSourceMinDiskSizeFieldName: int32(10),
			},
		},
	}
	flattened := flattenNotebookModels("test-depl-id", items.Items)
	assert.Equal(t, expected, flattened)
}

func TestFlattenNotebookModelList(t *testing.T) {
	items := []*nb.NotebookModel{
		{
			Id:          "test-list-id",
			Name:        "test-name",
			Cpu:         4,
			Memory:      10,
			MaxDiskSize: 10,
			MinDiskSize: 5,
		}}
	expected := []interface{}{
		map[string]interface{}{
			notebookModelDataSourceIdFieldName:          "test-list-id",
			notebookModelDataSourceNameFieldName:        "test-name",
			notebookModelDataSourceCPUFieldName:         float32(4),
			notebookModelDataSourceMemoryFieldName:      int32(10),
			notebookModelDataSourceMaxDiskSizeFieldName: int32(10),
			notebookModelDataSourceMinDiskSizeFieldName: int32(5),
		},
	}

	flattened := flattenNotebookModelList(items)
	assert.Equal(t, expected, flattened)
}
