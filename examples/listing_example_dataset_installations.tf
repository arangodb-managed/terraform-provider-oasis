/*
 * This example shows how to list all example dataset installations for a specific deployment.
 * It will fetch all installations based on the deployment id provided by the fetched deployment resource.
 * The data source will contain example installation id and the status of it.
 */
provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
  organization   = ""
}

// List all installations for a created deployment
data "oasis_example_dataset_installation" "my-installations" {
  deployment_id = oasis_deployment.my_oneshard_deployment.id
}

// We output the list of installations for this deployment
output "deployment-installations" {
  value = data.oasis_example_dataset_installation.my-installations
}

// Setup an oasis project
resource "oasis_project" "my_project" {
  name = "Test Terraform Project"
}

// Create / Read an oasis deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  name    = "Test Terraform Deployment"
  project = oasis_project.my_project.id
  location {
    region = "gcp-europe-west4"
  }
  version {
    // db_version = "3.6.0" // This is an optional field, if not set the default version will be used
  }
  configuration {
    model = "oneshard"
  }
}
