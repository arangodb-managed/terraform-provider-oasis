terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb.com/managed/oasis"
      version = ">=1.5.1"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
}

// Terraform created auditlog
resource "oasis_auditlog" "oasis_test_auditlog" {
  name         = "Terraform Oasis AuditLog"
  description  = "A test Oasis auditlog from Terraform Provider"
  organization = "" // organization id
}