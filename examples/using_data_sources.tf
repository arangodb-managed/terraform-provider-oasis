provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
}

resource "oasis_backup_policy" "my_backup_policy" {
  name = "Awesome Policy"
  description = "Description of the year"
  email_notification = "None"
  deployment_id = oasis_deployment.my_oneshard_deployment.id
  retention_period = 120
  upload = true
  schedule {
    type = "Monthly"
    monthly {
      day_of_month = 12
      schedule_at {
        hours = 15
        minutes = 10
        timezone = "UTC"
      }
    }
  }
}

resource "oasis_ipwhitelist" "my_iplist" {
  name = "terraform-ip-list"
  description = "Important ip list."
  cidr_ranges = ["1.2.3.4/32", "111.11.0.0/16", "0.0.0.0/0"]
  project = data.oasis_project.my_project.id
}

resource "oasis_deployment" "my_oneshard_deployment" {
  name = "terraform-deployment"
  description = "Description of the deployment"
  project      = data.oasis_project.my_project.id
  location {
    region   = "gcp-europe-west4"
  }
  version {
    db_version     = "3.6.0"
  }
  configuration {
    model          = "oneshard"
    node_count = 3
  }
  security {
    ip_whitelist = oasis_ipwhitelist.my_iplist.id
    ca_certificate = oasis_certificate.my_oasis_cert.id
  }
}

resource "oasis_certificate" "my_oasis_cert" {
  name = "terraform-cert"
  description = "Description of the certificate"
  project      = data.oasis_project.my_project.id
}

data "oasis_organization" "support" {
  id = "organization id"
}

data "oasis_project" "my_project" {
  id = "project id"
}

