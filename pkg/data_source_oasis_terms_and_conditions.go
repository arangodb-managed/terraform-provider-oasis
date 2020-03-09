//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package pkg

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// t&c data source fields
	tcIDFieldName           = "id"
	tcCreatedAtFieldName    = "created_at"
	tcContentFieldName      = "content"
	tcOrganizationFieldName = "organization"
)

// dataSourceTermsAndConditions defines a T&C datasource terraform type.
func dataSourceTermsAndConditions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTermsAndConditionsRead,

		Schema: map[string]*schema.Schema{
			tcIDFieldName: {
				Type:     schema.TypeString,
				Optional: true, // if left out, the current is fetched
			},
			tcContentFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tcCreatedAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tcOrganizationFieldName: {
				Type:     schema.TypeString,
				Optional: true, // if given, overwrites the plugin level organization
			},
		},
	}
}

// dataSourceTermsAndConditionsRead reloads the resource object from the terraform store.
func dataSourceTermsAndConditionsRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	var (
		tc  *rm.TermsAndConditions
		err error
	)
	if v, ok := data.GetOk(tcIDFieldName); ok {
		tc, err = rmc.GetTermsAndConditions(client.ctxWithToken, &common.IDOptions{Id: v.(string)})
		if err != nil {
			return err
		}
	} else {
		var orgID string
		if v, ok := data.GetOk(tcOrganizationFieldName); ok {
			orgID = v.(string)
		}
		tc, err = rmc.GetCurrentTermsAndConditions(client.ctxWithToken, &common.IDOptions{Id: orgID})
		if err != nil {
			return err
		}
	}

	for k, v := range flattenTCObject(tc) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(tc.GetId())
	return nil
}

// flattenTCObject creates a map from an Oasis Terms and Condition object for easy digestion by the terraform.
func flattenTCObject(tc *rm.TermsAndConditions) map[string]interface{} {
	return map[string]interface{}{
		tcIDFieldName:        tc.GetId(),
		tcCreatedAtFieldName: tc.GetCreatedAt().String(),
		tcContentFieldName:   tc.GetContent(),
	}
}
