---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oasis_organization Resource - terraform-provider-oasis"
subcategory: ""
description: |-
  Oasis Organization Resource
---

# oasis_organization (Resource)

Oasis Organization Resource

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Organization Resource Organization Name field

### Optional

- `description` (String) Organization Resource Organization Description field
- `locked` (Boolean) Organization Resource Organization Lock field

### Read-Only

- `id` (String) The ID of this resource.


