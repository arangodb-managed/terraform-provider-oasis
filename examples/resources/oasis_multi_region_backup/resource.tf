terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.6"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
  organization   = "" // Your Oasis organization where you want to create the resources
}

// Create Project
resource "oasis_project" "oasis_test_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
}

// Create Deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name                          = "oasis_multi_region_deployment"
  location {
    region = "gcp-europe-west4"
  }
  security {
    disable_foxx_authentication = false
  }
  disk_performance = "dp30"
  configuration {
    model                  = "oneshard"
    node_size_id           = "c4-a8"
    node_disk_size         = 20
    maximum_node_disk_size = 40
  }
  notification_settings {
    email_addresses = [
      "test@arangodb.com"
    ]
  }
}

// Create Backup
resource "oasis_backup" "backup" {
  name            = "oasis_backup"
  description     = "test backup description update from terraform"
  deployment_id   = oasis_deployment.my_oneshard_deployment.id
  upload          = true
  auto_deleted_at = 3 // auto delete after 3 days
}

// Create Multi Region Backup
resource "oasis_multi_region_backup" "backup" {
  source_backup_id = oasis_backup.backup.id // Existing backup ID
  region_id        = "gcp-us-central1"      // Oasis region identifier, which is other than the deployment region
}