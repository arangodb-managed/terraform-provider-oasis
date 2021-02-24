/*
 * This example demonstrates what fields are the bare minimum in order to get
 * started creating a deployment with Oasis.
 * The api_key_id and api_key_secret could be replaced with the environment
 * properties OASIS_API_KEY_ID and OASIS_API_KEY_SECRET respectively.
 */
provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
}

// A project is created here, but we could also use the default project
// which is created with a first organization.
resource "oasis_project" "my_project" {
  name = "Test Terraform Project"
}

/* A deployment has a bare minimum requirement of name, project, region,
 * version and model. Anything else is either left empty, or deciphered
 * for you. After a resource is created, `terraform show` can be used to
 * investiage calculated fields.
 * For example:
 * Calculated fields include, a certificate, node size, node disk size,
 * node memory and node count.
 */
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
