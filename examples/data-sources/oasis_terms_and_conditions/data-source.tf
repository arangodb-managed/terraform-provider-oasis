terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.0"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
}

// Load in Oasis terms and conditions
data "oasis_terms_and_conditions" "test_terms_and_conditions" {
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_terms_and_conditions.test_terms_and_conditions
}