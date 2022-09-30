terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.7"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
}

// Load in all the available datasets
data "oasis_notebook_model" "models" {
  deployment_id = "" // deployment id (required)
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_notebook_model.models
}