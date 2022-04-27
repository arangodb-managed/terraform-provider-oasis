terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source = "arangodb.com/managed/oasis"
      version = "1.5.1"
    }
  }
}

provider "oasis" {
  api_key_id     = "bqtrx1j8aoyjybqcpdw4" // API Key ID generated in Oasis platform
  api_key_secret = "44a268d1-50fd-27be-eac6-a81d6cfdf5b8" // API Key Secret generated in Oasis platform
}

// Terraform created organization
resource "oasis_organization" "oasis_test_organization" {
  name        = "Terraform Oasis Organization"
  description = "A test Oasis organization within from Terraform Provider"
}