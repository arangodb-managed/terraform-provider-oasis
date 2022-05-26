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

// Load in Oasis Cloud Providers
data "oasis_region" "test_oasis_region" {
  organization = "" // put your organization id here
  provider_id  = "" // put one of the cloud provider ids here (can be fetched from cloud_provider data source)
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_region.test_oasis_region
}