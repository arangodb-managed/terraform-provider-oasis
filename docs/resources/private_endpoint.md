---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oasis_private_endpoint Resource - terraform-provider-oasis"
subcategory: ""
description: |-
  Oasis Private Endpoint Resource
---

# oasis_private_endpoint (Resource)

Oasis Private Endpoint Resource

## Example Usage

```terraform
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.1"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
  organization   = ""
}

// Terraform created project
resource "oasis_project" "oasis_test_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
}

// Example of a oneshard deployment
resource "oasis_deployment" "my_aks_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name                          = "oasis_test_aks_dep_tf"

  location {
    region = "aks-westus2"
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

// Example of an AKS Private Endpoint
resource "oasis_private_endpoint" "my_aks_private_endpoint" {
  name        = "tf-private-endpoint-test"
  description = "Terraform generated AKS private endpoint"
  deployment  = oasis_deployment.my_aks_oneshard_deployment.id
  dns_names   = ["test.example.com", "test2.example.com"]
  aks {
    az_client_subscription_ids = ["291bba3f-e0a5-47bc-a099-3bdcb2a50a05"]
  }
}

// Example of an AWS oneshard deployment
resource "oasis_deployment" "my_aws_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name                          = "oasis_test_aws_dep_tf"

  location {
    region = "aws-us-east-2"
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

// Example of an AWS Private Endpoint
resource "oasis_private_endpoint" "my_aws_private_endpoint" {
  name        = "tf-private-endpoint-test"
  description = "Terraform generated AWS private endpoint"
  deployment  = oasis_deployment.my_aws_oneshard_deployment.id
  enable_private_dns = true
  dns_names   = ["test.example.com", "test2.example.com"]
  aws {
    principal {
      account_id = "123123123123"        // 12 digit AWS Account Identifier
      user_names = ["test@arangodb.com"] // User names (IAM User(s) that are able to setup the private endpoint)
      role_names = ["test"]              // Role names (IAM role(s) that are able to setup the endpoint)
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `deployment` (String) Private Endpoint Resource Private Endpoint Deployment ID field
- `name` (String) Private Endpoint Resource Private Endpoint Name field

### Optional

- `aks` (Block List, Max: 1) Private Endpoint Resource Private Endpoint AKS field (see [below for nested schema](#nestedblock--aks))
- `aws` (Block List, Max: 1) Private Endpoint Resource Private Endpoint AWS field (see [below for nested schema](#nestedblock--aws))
- `description` (String) Private Endpoint Resource Private Endpoint Description field
- `enable_private_dns` (Bool) If set, private DNS zone integration is enabled for this private endpoint service. For GCP this bool is immutable, so can only be set during the creation. For AKS this boolean cannot be set.
- `dns_names` (List of String) Private Endpoint Resource Private Endpoint DNS Names field (list of dns names)
- `gcp` (Block List, Max: 1) Private Endpoint Resource Private Endpoint GCP field (see [below for nested schema](#nestedblock--gcp))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--aks"></a>
### Nested Schema for `aks`

Required:

- `az_client_subscription_ids` (List of String) Private Endpoint Resource Private Endpoint AKS Subscription IDS field (list of subscription ids)


<a id="nestedblock--aws"></a>
### Nested Schema for `aws`

Required:

- `principal` (Block List, Min: 1) Private Endpoint Resource Private Endpoint AWS Principal field (see [below for nested schema](#nestedblock--aws--principal))

<a id="nestedblock--aws--principal"></a>
### Nested Schema for `aws.principal`

Required:

- `account_id` (String) Private Endpoint Resource Private Endpoint AWS Principal Account Id field

Optional:

- `role_names` (List of String) Private Endpoint Resource Private Endpoint AWS Principal Role Names field
- `user_names` (List of String) Private Endpoint Resource Private Endpoint AWS Principal User Names field



<a id="nestedblock--gcp"></a>
### Nested Schema for `gcp`

Required:

- `projects` (List of String) Private Endpoint Resource Private Endpoint GCP Projects field (list of project ids)


