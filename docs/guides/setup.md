---
page_title: "Setup"
description: |-
    Guide to setup the Oasis Provider with the required Keys
---

# Setup

When using the provider, you need to set up api_key_id and api_key_secret. Both parameters can be generated in the Oasis Dashboard. Once you are logged in, navigate to the [**API Keys**](https://cloud.arangodb.com/dashboard/user/api-keys) tab of your user account and click the **New API** key button.

```hcl
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
```

The provider can also be setup with a default organization and project to manage resources in:

```hcl
provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
  organization   = "" // Oasis Organization ID
  project        = "" // Project ID within the specified organization
}
```

The other options you can provide are:
- `oasis_endpoint` for the endpoint you want to manage the resources in (by default set to: `api.cloud.arangodb.com`).
- `api_port_suffix` for the Oasis API Port Suffix (by default set to `:443`).