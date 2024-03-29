---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oasis_notebook Resource - terraform-provider-oasis"
subcategory: ""
description: |-
  Oasis Notebook Resource
---

# oasis_notebook (Resource)

Oasis Notebook Resource

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

// Create Project
resource "oasis_project" "oasis_test_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
}

// Create Deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project                       = oasis_project.oasis_test_project.id // Project id where deployment will be created
  name                          = "oasis_jupyter_notebook_deployment"
  location {
    region = "gcp-europe-west4"
  }
  security {
    disable_foxx_authentication = false
  }
  disk_performance = "dp30"
  configuration {
    model                  = "oneshard"
    node_size_id           = "c4-a8"
    node_disk_size         = 20
    maximum_node_disk_size = 40
  }
  notification_settings {
    email_addresses = [
      "test@arangodb.com"
    ]
  }
}

// Create Notebook
resource "oasis_notebook" "oasis_test_notebook" {
  deployment_id = oasis_deployment.my_oneshard_deployment.id
  name          = "Test Oasis Jupyter Notebook"
  description   = "Test Description"
  model {
    notebook_model_id = "basic"
    disk_size         = "10"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `deployment_id` (String) Notebook Resource Notebook Deployment ID field
- `model` (Block List, Min: 1, Max: 1) Notebook Resource Notebook Model field (see [below for nested schema](#nestedblock--model))
- `name` (String) Notebook Resource Notebook Name field

### Optional

- `description` (String) Notebook Resource Notebook Description field
- `url` (String) Notebook Resource Notebook URL field

### Read-Only

- `created_at` (String) Notebook Resource Notebook Created At field
- `created_by_id` (String) Notebook Resource Notebook Created By Id field
- `deleted_at` (String) Notebook Resource Notebook Deleted At field
- `id` (String) The ID of this resource.
- `is_deleted` (Boolean) Notebook Resource Notebook Is Deleted field
- `is_paused` (Boolean) Notebook Resource Notebook Is Paused field
- `last_paused_at` (String) Notebook Resource Notebook Last Paused field
- `last_resumed_at` (String) Notebook Resource Notebook Last Resumed field
- `status` (List of Object) Notebook Resource Notebook Status field (see [below for nested schema](#nestedatt--status))

<a id="nestedblock--model"></a>
### Nested Schema for `model`

Required:

- `disk_size` (Number) Notebook Resource Notebook Model Disk Size field
- `notebook_model_id` (String) Notebook Resource Notebook Model ID field


<a id="nestedatt--status"></a>
### Nested Schema for `status`

Read-Only:

- `endpoint` (String)
- `last_updated_at` (String)
- `message` (String)
- `phase` (String)
- `usage` (List of Object) (see [below for nested schema](#nestedobjatt--status--usage))

<a id="nestedobjatt--status--usage"></a>
### Nested Schema for `status.usage`

Read-Only:

- `last_cpu_limit` (Number)
- `last_cpu_usage` (Number)
- `last_memory_limit` (Number)
- `last_memory_usage` (Number)


