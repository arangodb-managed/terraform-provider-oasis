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
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
)

const (
	// current user source fields
	userIdFieldName    = "id"
	userEmailFieldName = "email"
	userNameFieldName  = "name"
)

// dataSourceOasisCurrentUser defines a Current User datasource terraform type.
func dataSourceOasisCurrentUser() *schema.Resource {
	return &schema.Resource{
		Description: "Oasis Current User Data Source",
		ReadContext: dataSourceOasisCurrentUserRead,

		Schema: map[string]*schema.Schema{
			userIdFieldName: {
				Type:        schema.TypeString,
				Description: "Current User Data Source User ID field",
				Optional:    true,
			},
			userEmailFieldName: {
				Type:        schema.TypeString,
				Description: "Current User Data Source Email field",
				Optional:    true,
			},
			userNameFieldName: {
				Type:        schema.TypeString,
				Description: "Current User Data Source Name field",
				Optional:    true,
			},
		},
	}
}

// dataSourceOasisCurrentUserRead reloads the resource object from the terraform store.
func dataSourceOasisCurrentUserRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return diag.FromErr(err)
	}

	iamc := iam.NewIAMServiceClient(client.conn)
	user, err := iamc.GetThisUser(client.ctxWithToken, &common.Empty{})
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range flattenCurrentUserObject(user) {
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	data.SetId(user.GetId())
	return nil
}

// flattenCurrentUserObject creates a map from an Oasis Current User for easy digestion by the terraform schema.
func flattenCurrentUserObject(user *iam.User) map[string]interface{} {
	return map[string]interface{}{
		userIdFieldName:    user.GetId(),
		userEmailFieldName: user.GetEmail(),
		userNameFieldName:  user.GetName(),
	}
}
