provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
  organization   = "organization id"
  # project = "1234567" -- this is optional
}

resource "oasis_project" "my_project" {
  name = "Test Terraform Project 2"
  description = "This should also be filled 1"
  # organization = "This is optional"
}

resource "oasis_deployment" "my_oneshard_deployment" {
  name = "terraform-deployment"
  project      = oasis_project.my_project.id
  location {
    region   = "gcp-europe-west4"
  }
  version {
    db_version     = "3.6.0"
  }
  # Model is oneshard and the count is 3 by default.
  configuration {
    model = "oneshard" # this is a required field
  #     node_count = 3
  #   }
  #   security {
  #     ca_certificate = oasis_certificate.my_oasis_cert.id # if non is provided, one will be created and automatically used.
  }
}
