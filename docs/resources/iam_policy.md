---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oasis_iam_policy Resource - terraform-provider-oasis"
subcategory: ""
description: |-
  Oasis IAM Policy Resource
---

# oasis_iam_policy (Resource)

Oasis IAM Policy Resource

## Example Usage

```terraform
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.2"
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

// Load in an Oasis Current User within an organization
data "oasis_current_user" "oasis_test_current_user" {}

// Terraform created IAM Policy. This resource uses the computed ID value of the
// previously defined organization resource and IAM group resource.
resource "oasis_iam_policy" "my_iam_policy_group" {
  url = "/Organization/${oasis_organization.oasis_test_organization.id}"

  binding {
    role  = "auditlog-admin"
    group = oasis_iam_group.my_iam_group.id
  }

  binding {
    role  = "auditlog-archive-viewer"
    group = oasis_iam_group.my_iam_group.id
  }

  binding {
    role = "auditlog-archive-viewer"
    user = data.oasis_current_user.oasis_test_current_user.id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `binding` (Block List, Min: 1) IAM Policy Resource IAM Policy Bindings (see [below for nested schema](#nestedblock--binding))
- `url` (String) IAM Policy Resource IAM Policy URL

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--binding"></a>
### Nested Schema for `binding`

Required:

- `role` (String) IAM Policy Resource IAM Policy Role

Optional:

- `group` (String) IAM Policy Resource IAM Policy Group
- `user` (String) IAM Policy Resource IAM Policy User


