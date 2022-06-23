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
  organization   = "" // Your Oasis organization where you want to create the resources
}

// Load in an Oasis backup of a deployment
data "oasis_backup" "oasis_test_backup" {
  id = "" // Provide your Backup ID here
}

// Output the data after it has been synced.
output "backup" {
  value = data.oasis_backup.oasis_test_backup
}