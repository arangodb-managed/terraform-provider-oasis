terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.0"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
  organization   = "" // Your Oasis organization where you want to create the resources
}

// Terraform created project.
resource "oasis_project" "my_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
}

// Terraform created ip allowlist. This resource uses the computed ID value of the
// previously defined project resource.
resource "oasis_ipallowlist" "my_iplist" {
  name        = "Terraform IP Allowlist"
  description = "IP Allowlist description"
  cidr_ranges = ["1.2.3.4/32", "111.11.0.0/16"]
  project     = oasis_project.my_project.id
}