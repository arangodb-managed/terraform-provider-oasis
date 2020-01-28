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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

const (
	// Organization data source fields
	idFieldName                             = "id"
	nameFieldName                           = "name"
	descriptionFieldName                    = "description"
	urlFieldName                            = "url"
	createdAtFieldName                      = "created_at"
	isDeletedFieldName                      = "is_deleted"
	tierFieldName                           = "tier"
	tierIdFieldName                         = "id"
	tierNameFieldName                       = "name"
	tierHasSupportPlansFieldName            = "has_support_plans"
	tierHasBackupUploadsFieldName           = "has_backup_uploads"
	tierRequiresTermsAndConditionsFieldName = "requires_terms_and_conditions"
)

// dataSourceOasisOrganization defines an Organization datasource terraform type.
func dataSourceOasisOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOasisOrganizationRead,

		Schema: map[string]*schema.Schema{
			idFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			nameFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			descriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			urlFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			createdAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			isDeletedFieldName: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			tierFieldName: {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tierIdFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tierNameFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tierHasSupportPlansFieldName: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						tierHasBackupUploadsFieldName: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						tierRequiresTermsAndConditionsFieldName: {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceOasisOrganizationRead reloads the resource object from the terraform store.
func dataSourceOasisOrganizationRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return err
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	oid := data.Get(idFieldName).(string)
	org, err := rmc.GetOrganization(client.ctxWithToken, &common.IDOptions{Id: oid})
	if err != nil {
		return err
	}
	for k, v := range flattenOrganizationObject(org) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(org.GetId())
	return nil
}

// flattenOrganizationObject creates a map from an Oasis Organization structure for the terraform schema.
func flattenOrganizationObject(org *rm.Organization) map[string]interface{} {
	ret := map[string]interface{}{
		idFieldName:          org.GetId(),
		nameFieldName:        org.GetName(),
		descriptionFieldName: org.GetDescription(),
		urlFieldName:         org.GetUrl(),
		createdAtFieldName:   org.GetCreatedAt().String(),
		tierFieldName:        flattenTierObject(org.GetTier()),
	}

	return ret
}

// flattenTierObject will produce an accepted transformation of the tier struct to terraform schema.
// Note that boolean values are transformed to string, because at the time of this, 0.12 terraform
// does not support different value types inside a map structure.
func flattenTierObject(tier *rm.Tier) interface{} {
	return map[string]interface{}{
		tierIdFieldName:                         tier.GetId(),
		tierNameFieldName:                       tier.GetName(),
		tierHasSupportPlansFieldName:            fmt.Sprintf("%t", tier.GetHasSupportPlans()),
		tierHasBackupUploadsFieldName:           fmt.Sprintf("%t", tier.GetHasBackupUploads()),
		tierRequiresTermsAndConditionsFieldName: fmt.Sprintf("%t", tier.GetRequiresTermsAndConditions()),
	}
}
