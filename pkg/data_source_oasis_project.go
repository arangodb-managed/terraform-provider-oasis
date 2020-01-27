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

func dataSourceOasisProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOasisProjectRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOasisProjectRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.Connect()
	if err != nil {
		return err
	}

	rmc := rm.NewResourceManagerServiceClient(client.conn)
	id := data.Get("id").(string)
	proj, err := rmc.GetProject(client.ctxWithToken, &common.IDOptions{Id: id})
	if err != nil {
		return err
	}

	for k, v := range flattenProjectObject(proj) {
		err := data.Set(k, v)
		if err != nil {
			return err
		}
	}
	data.SetId(proj.GetId())
	return nil
}

func flattenProjectObject(proj *rm.Project) map[string]interface{} {
	return map[string]interface{}{
		"id":          proj.GetId(),
		"name":        proj.GetName(),
		"description": proj.GetDescription(),
		"url":         proj.GetUrl(),
		"created_at":  proj.GetCreatedAt().String(),
	}
}
