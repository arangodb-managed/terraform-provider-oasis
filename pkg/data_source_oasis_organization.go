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
	// Organization data source fields
	orgIdFieldName                          = "id"
	orgNameFieldName                        = "name"
	orgDescriptionFieldName                 = "description"
	orgUrlFieldName                         = "url"
	orgCreatedAtFieldName                   = "created_at"
	orgIsDeletedFieldName                   = "is_deleted"
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
			orgIdFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			orgNameFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			orgDescriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			orgUrlFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			orgCreatedAtFieldName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			orgIsDeletedFieldName: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			tierFieldName: {
				Type:     schema.TypeSet,
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
	oid := data.Get(orgIdFieldName).(string)
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
		orgIdFieldName:          org.GetId(),
		orgNameFieldName:        org.GetName(),
		orgDescriptionFieldName: org.GetDescription(),
		orgUrlFieldName:         org.GetUrl(),
		orgCreatedAtFieldName:   org.GetCreatedAt().String(),
		tierFieldName:           flattenTierObject(org.GetTier()),
	}

	return ret
}

// flattenTierObject will produce a schema set which can be interpreted by terraform for type safety.
// Contrary to first look, tier is not handled as TypeMap because in terraform 0.12, there is no support
// for different embedded types. The hash for the set is based on the definition of the schema for tier.
func flattenTierObject(tier *rm.Tier) *schema.Set {
	s := &schema.Set{
		F: schema.HashResource(dataSourceOasisOrganization().Schema[tierFieldName].Elem.(*schema.Resource)),
	}
	tierMap := map[string]interface{}{
		tierIdFieldName:                         tier.GetId(),
		tierNameFieldName:                       tier.GetName(),
		tierHasSupportPlansFieldName:            tier.GetHasSupportPlans(),
		tierHasBackupUploadsFieldName:           tier.GetHasBackupUploads(),
		tierRequiresTermsAndConditionsFieldName: tier.GetRequiresTermsAndConditions(),
	}
	s.Add(tierMap)
	return s
}
