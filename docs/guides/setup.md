---
page_title: "Setup"
description: |-
    Guide to setup the Oasis Provider with the required Keys
---

# Setup

When using the provider you need to setup `api_key_id` and `api_key_secret`. The two API keys can be generated from Oasis Dashboard. On a logged in view, navigate to [API Keys](https://cloud.arangodb.com/dashboard/user/api-keys)

```terraform
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
  organization   = "" // Your Oasis organization where you want to create the resources
}
```

If you already followed the setup instructions from Getting Started guide, you can use those keys here.