terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source = "arangodb.com/managed/oasis"
      version = ">=1.5.1"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
  organization   = "" // Your Oasis organization where you want to create the resources
}

// Terraform created project.
resource "oasis_project" "oasis_test_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
}

// Example of a oneshard deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name = "oasis_test_dep_tf"
  location {
    region = "gcp-europe-west4"
  }
  version {
    db_version = "3.8.6"
  }
  configuration {
    model = "oneshard"
    node_size_id = "c4-a8"
    node_disk_size = 20
  }
  notification_settings {
    email_addresses = [
      "test@arangodb.com"
    ]
  }
}

// Oasis backup
// This resources uses the computed ID of the deployment created above.
resource "oasis_backup" "my_backup" {
  name = "test tf backup"
  description = "test backup description from terraform"
  deployment_id = oasis_deployment.my_oneshard_deployment.id
  upload = true
  auto_deleted_at = -3 // auto delete after 3 days
}