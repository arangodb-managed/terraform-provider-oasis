terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source = "arangodb.com/managed/oasis"
      version = "1.5.1"
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

// Example of a one-shard deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project = oasis_project.oasis_test_project.id // If set here, overrides project in provider
  name = "oasis_test_dep_tf"
  location {
    region = "gcp-europe-west4" // Required
  }
  version { // this section is optional
    db_version = "3.8.6" // This is an optional field, if not set the default version will be used
  }
  security { // this section is optional
    ca_certificate = "" // If not set, uses default certificate from project
    ip_allowlist = "" // If not set, no allowlist is configured
    disable_foxx_authentication = false // If set to true, request to Foxx apps are not authentications.
  }
  configuration {
    model = "oneshard"
    node_size_id = "a4"
    node_disk_size = 20
  }
  notification_settings {
    email_addresses = [
      "test@arangodb.com" // this will set email addresses used for notifications regarding depoyment
    ]
  }
}

// Oasis backup policy. This can have a lot of values and configuration options.
// For details, please consult `terraform providers schema` or the code.
// This resources uses the computed ID of the deployment created above.
resource "oasis_backup_policy" "my_backup_policy" {
  name               = "Test Policy"
  description        = "Test Description"
  email_notification = "FailureOnly"
  deployment_id      = oasis_deployment.my_oneshard_deployment.id
  retention_period_hour   = 120
  upload             = true
  schedule {
    type = "Monthly"
    monthly {
      day_of_month = 12
      schedule_at {
        hours    = 15
        minutes  = 10
        timezone = "UTC"
      }
    }
  }
}