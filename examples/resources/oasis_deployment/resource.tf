terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.8"
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
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name                          = "oasis_test_dep_tf"

  location {
    region = "gcp-europe-west4"
  }

  configuration {
    model = "oneshard"
  }

  notification_settings {
    email_addresses = [
      "test@arangodb.com"
    ]
  }
}


// Example of a sharded deployment
resource "oasis_deployment" "my_sharded_deployment" {
  terms_and_conditions_accepted = "true"
  name                          = "oasis_sharded_dep_tf"
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  location {
    region = "gcp-europe-west4"
  }

  version {
    db_version = "3.9.1"
  }

  security {                            // this section is optional
    ca_certificate              = ""    // If not set, uses default certificate from project (this is here as an empty string for documentation purposes)
    ip_allowlist                = ""    // If not set, no allowlist is configured (this is here as an empty string for documentation purposes)
    disable_foxx_authentication = false // If set to true, request to Foxx apps are not authentications.
  }

  configuration {
    model          = "sharded"
    node_size_id   = "c4-a4"
    node_disk_size = 20
    node_count     = 5
  }
}