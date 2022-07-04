terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.2"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
}

// Load in an Oasis Current User within an organization
data "oasis_current_user" "oasis_test_current_user" {}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_current_user.oasis_test_current_user
}