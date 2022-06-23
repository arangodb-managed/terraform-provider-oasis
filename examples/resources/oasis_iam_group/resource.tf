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
}

// Terraform created organization
resource "oasis_organization" "oasis_test_organization" {
  name        = "Terraform Oasis Organization"
  description = "A test Oasis organization from Terraform Provider"
}

// Terraform created IAM Group. This resource uses the computed ID value of the
// previously defined organization resource.
resource "oasis_iam_group" "my_iam_group" {
  name         = "Terraform IAM Group"
  description  = "IAM Group created by Terraform"
  organization = oasis_organization.oasis_test_organization.id
}