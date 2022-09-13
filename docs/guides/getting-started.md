---
page_title: "Getting started with Oasis Terraform Provider"
description: |-
    Guide to getting started with the ArangoDB Cloud (Oasis) Terraform Provider
---

# Getting started with Oasis Terraform Provider
ArangoDB Cloud (Oasis), provides ArangoDB databases as a Service (DBaaS). It enables you to use the entire functionality of an ArangoDB deployment without the need to run or manage the system yourself.

Terraform Provider Oasis is a plugin for Terraform that allows for the full lifecycle management of ArangoDB Cloud (Oasis) resources.

## Provider Setup


You need to supply proper credentials to the provider before it can be used. API keys serve as the credentials to the provider. You can obtain the keys from the Oasis dashboard.
Log in to the Oasis dashboard and open the [**API Keys**](https://cloud.arangodb.com/dashboard/user/api-keys) tab of your user account. Click the **New API key** button to generate a new key, which can be used with ArangoDB's public API.

```hcl
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
```
Optionally, you can provide `oasis_endpoint`, `api_port_suffix`, `organization` or `project`.  Please refer to the [Setup](setup.md) section for more information.

## Example Usage

Creating your first Deployment:

```hcl
# Terraform created organization
resource "oasis_organization" "oasis_test_organization" {
  name        = "Terraform Oasis Organization"
  description = "A test Oasis organization from Terraform Provider"
}

# Terraform created project.
resource "oasis_project" "oasis_test_project" {
  name         = "Terraform Oasis Project"
  description  = "A test Oasis project within an organization from the Terraform Provider"
  organization = oasis_organization.oasis_test_organization.id
}

# Example of a oneshard deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name                          = "oasis_test_dep_tf"

  location {
    region = "gcp-europe-west4"
  }

  version {
    db_version = "3.8.7"
  }

  configuration {
    model = "oneshard"
  }

  notification_settings {
    email_addresses = [
      "test@arangodb.com"
    ]
  }
}
```