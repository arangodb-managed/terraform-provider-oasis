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


resource "oasis_multi_region_backup" "backup" {
  source_backup_id = "" // Existing backup ID that is already uploaded
  region_id        = "" // Oasis region identifier, which is other than the deployment region
}