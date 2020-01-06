ArangoDB Oasis Terraform Provider
==================

<img src="https://cloud.arangodb.com/static/media/cloud.5973146f.svg" width="300px">

- Website: https://cloud.arangodb.com/
- Slack: https://slack.arangodb.com/


Maintainers
-----------

This provider plugin is maintained by the team at [ArangoDB](https://www.arangodb.com/).


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Usage
---------------------
TODOs:
Explain how to:
* Create/Retrieve API Key
* Retrieve relevant ids (org, project, provider, ca_certificate, ip_whitelist)

```
provider "oasis" {
  api_key_id = "xx"
  api_key_secret  = "xxx"
}

resource "oasis_deployment" "my_deployment" {
  organization = "190765105"

  project = "190765139"

  location = {
    provider = ""
    region = ""
  }

  version = {
    db_version = ""
    ca_certificate = ""
    ip_whitelist = ""
  }

  configuration = {
    sharded =  false
    node_memory_gb = 4
    node_disk_gb = 10
    num_nodes = 3
  }
}


```


