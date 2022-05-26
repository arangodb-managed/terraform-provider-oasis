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

// Load in an Oasis project within an organization
data "oasis_project" "oasis_test_project" {
  id = "" // Provide your project ID here
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_project.oasis_test_project
}