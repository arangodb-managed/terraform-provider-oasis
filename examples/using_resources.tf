/*
 * This example shows how to use terraform to create everything
 * except for the organization. But even that could be just a data source
 * or an environment property.
 */
provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
  organization   = "_support"
}

// Terraform created project.
resource "oasis_project" "my_project" {
  name        = "Terraform Project"
  description = "Project description"
}

// Terraform created ip whitelist. This resource uses the computed ID value of the
// previously defined project resource.
resource "oasis_ipwhitelist" "my_iplist" {
  name        = "Terraform IP Whitelist"
  description = "IP Whitelist description"
  cidr_ranges = ["1.2.3.4/32", "111.11.0.0/16", "0.0.0.0/0"]
  project     = oasis_project.my_project.id
}

// Terraform created deployment. For all fields, please consult `terraform providers schema`
// or the code.
// This resource uses the computed project ID of the previously defined project resource,
// and two other resources, ip whitelist and the certificate.
resource "oasis_deployment" "my_oneshard_deployment" {
  name        = "Terraform Deployment"
  description = "Deployment description"
  project     = oasis_project.my_project.id
  location {
    region = "gcp-europe-west4"
  }
  version {
    db_version = "3.6.0"
  }
  configuration {
    model      = "oneshard"
    node_count = 3
  }
  security {
    ip_whitelist   = oasis_ipwhitelist.my_iplist.id
    ca_certificate = oasis_certificate.my_oasis_cert.id
  }
}

// Oasis backup policy. This can have a lot of values and configuration options.
// For details, please consult `terraform providers schema` or the code.
// This resources uses the computed ID of the deployment created above.
resource "oasis_backup_policy" "my_backup_policy" {
  name               = "Awesome Policy"
  description        = "Description of the year"
  email_notification = "None"
  deployment_id      = oasis_deployment.my_oneshard_deployment.id
  retention_period   = 120
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

// Oasis certificate. This certificate's computed ID will be used in the project above.
resource "oasis_certificate" "my_oasis_cert" {
  name        = "Terraform certificate"
  description = "Certificate description."
  project     = oasis_project.my_project.id
}
