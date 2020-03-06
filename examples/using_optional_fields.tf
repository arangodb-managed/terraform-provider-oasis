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
  configuration {
    model = "oneshard" # this is a required field
    # this is an optional field and automatically set.
    # further more, the smallest node size available in the given region will be used.
    #node_count = 3
  }
  # Security configuration is optional.
  # If no certificate is provided, one will be generated or the default will be used.
  #security {
  #  ca_certificate = oasis_certificate.my_oasis_cert.id
  #}
}
