---
page_title: "Getting started with Oasis Terraform Provider"
description: |-
    Guide to getting started with the ArangoDB Oasis Terraform Provider
---

# Getting started with Oasis Terraform Provider
ArangoDB Oasis, the ArangoDB Cloud, provides ArangoDB databases as a Service (DBaaS). It enables you to use the entire functionality of an ArangoDB cluster deployment without the need to run or manage the system yourself.

Terraform Provider Oasis is a plugin for Terraform that allows for the full lifecycle management of ArangoDB Cloud resources.

## Provider Setup

The provider needs to be configured with the proper credentials before it can be used. You will need API Keys to interact with the Terraform Provider. Api Keys can be generated and viewed under the user's dashboard view on the API Keys tab.
On a logged in view, navigate to [API Keys](https://cloud.arangodb.com/dashboard/user/api-keys) and hit the button
labeled `New API key`. This will generate a set of keys which can be used with ArangoDB's public API.

```hcl
provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
}
```
Optionally, you can provide `oasis_endpoint`, `api_port_suffix`, `organization` or `project`. Please visit the provider initialization documentation for more information.

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
    db_version = "3.8.6"
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