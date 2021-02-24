/*
 * This example will demonstrate how to install an example dataset for a finished deployment.
 * The below example dataset installation resource will take the dataset with id imdb
 * and import it into the deployment once it finishes bootstrapping. After that, the data
 * will be loaded in a random generated db name ( which will be displayed in the status output ).
 */
provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
}

// Use an existing project
data "oasis_project" "my_project" {
  id = "" // enter existing project ID here
}

/*
 * Create a deployment.
 */
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  name    = "Test Terraform Deployment"
  project = data.oasis_project.my_project.id
  location {
    region = "gcp-europe-west4"
  }
  version {
    // db_version = "3.6.2.2" // This is an optional field, if not set the default version will be used
  }
  configuration {
    model = "oneshard"
  }
}

/*
 * Create an example dataset installation for the deployment to have some data
 * to play with once it finishes bootstrapping.
 */
resource "oasis_example_dataset_installation" "imdb-movie-data" {
  deployment_id = oasis_deployment.my_oneshard_deployment.id
  example_dataset_id = "imdb"
}