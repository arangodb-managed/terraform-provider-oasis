terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb.com/managed/oasis"
      version = ">=1.5.1"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
}

// Load in an Oasis organization
data "oasis_organization" "oasis_test_organization" {
  id = "" // Provide your Organization ID here
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_organization.oasis_test_organization
}