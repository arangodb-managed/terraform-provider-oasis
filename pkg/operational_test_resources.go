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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// FetchOrganizationID finds and retrieves the first Organization ID it finds for a user.
func FetchOrganizationID(testAccProvider *schema.Provider) (string, error) {
	// Initialize Client with connection settings
	if err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil)); err != nil {
		return "", err
	}
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return "", err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	if organizations, err := rmc.ListOrganizations(client.ctxWithToken, &common.ListOptions{}); err != nil {
		client.log.Error().Err(err).Msg("Failed to list organizations")
		return "", err
	} else if len(organizations.GetItems()) < 1 {
		client.log.Error().Err(err).Msg("No organizations found")
		return "", fmt.Errorf("no organizations found")
	} else {
		return organizations.GetItems()[0].GetId(), nil
	}
}

// FetchProjectID will find the first project given an organization and retrieve its ID.
func FetchProjectID(orgID string, testAccProvider *schema.Provider) (string, error) {
	// Initialize Client with connection settings
	if err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil)); err != nil {
		return "", err
	}
	client := testAccProvider.Meta().(*Client)
	if err := client.Connect(); err != nil {
		client.log.Error().Err(err).Msg("Failed to connect to api")
		return "", err
	}
	rmc := rm.NewResourceManagerServiceClient(client.conn)
	if projects, err := rmc.ListProjects(client.ctxWithToken, &common.ListOptions{ContextId: orgID}); err != nil {
		client.log.Error().Err(err).Str("organization-id", orgID).Msg("Failed to list projects for organization")
		return "", err
	} else if len(projects.GetItems()) < 1 {
		client.log.Error().Err(err).Str("organization-id", orgID).Msg("No projects found")
		return "", fmt.Errorf("no projects found")
	} else {
		return projects.GetItems()[0].GetId(), nil
	}
}