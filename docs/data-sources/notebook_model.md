---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oasis_notebook_model Data Source - terraform-provider-oasis"
subcategory: ""
description: |-
  Oasis Notebook Model Data Source
---

# oasis_notebook_model (Data Source)

Oasis Notebook Model Data Source

## Example Usage

```terraform
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.7"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
}

// Load in all the available datasets
data "oasis_notebook_model" "models" {
  deployment_id = "" // deployment id (required)
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_notebook_model.models
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `deployment_id` (String) Notebook Model Data Source Notebook Model Deployment ID field

### Read-Only

- `id` (String) The ID of this resource.
- `items` (List of Object) (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `cpu` (Number)
- `id` (String)
- `max_disk_size` (Number)
- `memory` (Number)
- `min_disk_size` (Number)
- `name` (String)

