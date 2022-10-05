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

// Create Project
resource "oasis_project" "oasis_test_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
}

// Create Deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name                          = "oasis_jupyter_notebook_deployment"
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

// Create Notebook
resource "oasis_notebook" "oasis_test_notebook" {
  deployment_id = oasis_deployment.my_oneshard_deployment.id
  name          = "Test Oasis Jupyter Notebook"
  description   = "Test Description"
  model {
    notebook_model_id = "basic"
    disk_size         = "10"
  }
}