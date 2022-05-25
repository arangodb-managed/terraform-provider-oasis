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
  organization   = ""
}

// Terraform created project
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
    region = "aks-westus2"
  }

  version {
    db_version = "3.8.6"
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

// Example of a Private Endpoint
resource "oasis_private_endpoint" "my_private_endpoint" {
  name        = "tf-private-endpoint-test"
  description = "Terraform generated private endpoint"
  deployment  = oasis_deployment.my_oneshard_deployment.id
  dns_names   = ["test.example.com", "test2.example.com"]
  aks {
    az_client_subscription_ids = ["291bba3f-e0a5-47bc-a099-3bdcb2a50a05"]
  }
}