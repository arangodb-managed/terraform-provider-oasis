/*
 * This example shows how to use data sources in oasis.
 * Specifically, using organizations and projects as data sources.
 * If there are already projects which should be used with the generated
 * resources, define a data source with the ID of the project or the
 * organization to use.
 *
 * Data sources can be used with the `data.` prefix.
 * The definitions can be in any order. Terraform assures that all dependencies
 * will be resolved first.
 */
provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
}

// Define existing organizations with oasis_organization.
data "oasis_organization" "support" {
  id = "organization id"
}

// Define existing projects with oasis_project.
data "oasis_project" "my_project" {
  id = "project id"
}

// Backup policy definition. This requires a deployment_id which is provided by the created deployment.
// Backup policy has a lot of configuration options variations for timing.
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

// IP whitelist. This needs a project field.
resource "oasis_ipwhitelist" "my_iplist" {
  name        = "terraform-ip-list"
  description = "Important ip list."
  cidr_ranges = ["1.2.3.4/32", "111.11.0.0/16", "0.0.0.0/0"]
  project     = data.oasis_project.my_project.id
}

// Deployment has a lot of moving parts. There are not that many though which are requires fields.
// This deployment is of type oneshard and defines 3 database nodes. Though note that the node_count
// field is actually optional.
resource "oasis_deployment" "my_oneshard_deployment" {
  name        = "terraform-deployment"
  description = "Description of the deployment"
  project     = data.oasis_project.my_project.id
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

// Create a certificate. This is also an optional thing to do, because a new, default certificate is
// created and provided with all new deployments. We just create a specific one here with a specific name.
resource "oasis_certificate" "my_oasis_cert" {
  name        = "terraform-cert"
  description = "Description of the certificate"
  project     = data.oasis_project.my_project.id
}
