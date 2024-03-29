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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/assert"

	iam "github.com/arangodb-managed/apis/iam/v1"
)

func TestFlattenCurrentUserDataSource(t *testing.T) {
	testUserId := fmt.Sprintf("-oauth2|%s", acctest.RandString(21))
	user := iam.User{
		Id:    testUserId,
		Name:  "Test",
		Email: "test@arangodb.com",
	}
	expected := map[string]interface{}{
		userIdFieldName:    testUserId,
		userNameFieldName:  "Test",
		userEmailFieldName: "test@arangodb.com",
	}
	got := flattenCurrentUserObject(&user)
	assert.Equal(t, expected, got)
}
