//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
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

func TestFlattenExampleDatasetInstallations(t *testing.T) {
	created, _ := types.TimestampProto(time.Date(1980, 03, 03, 1, 1, 1, 0, time.UTC))
	items := &example.ExampleDatasetInstallationList{
		Items: []*example.ExampleDatasetInstallation{
			{
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
			},
		},
	}
	expected := map[string]interface{}{
		installationDeploymentIdFieldName: "test-depl-id",
		installationItemsFieldName: []interface{}{
			map[string]interface{}{
				installationExampleDatasetIdFieldName: "test-example-id",
				installationCreatedAtFieldName:        "1980-03-03T01:01:01Z",
				installationStatusFieldName: []interface{}{
					map[string]interface{}{
						installationStatusStateFieldName:        example.ExampleInstallationStatusReady,
						installationStatusIsFailedFieldName:     false,
						installationStatusIsAvailableFieldName:  true,
						installationStatusDatabaseNameFieldName: "test-db",
					},
				},
			},
		},
	}
	flattened := flattenExampleDatasetInstallations("test-depl-id", items.Items)
	assert.Equal(t, expected, flattened)
}
