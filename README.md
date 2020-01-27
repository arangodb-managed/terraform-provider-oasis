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
  organization = "190765105"
  project = "foo"
}

// Example of oneshard deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  project = "190765139" // If set here, overrides project in provider
  location = {
    provider = "gcp" // Required
    region = "gcp-europe-west4" // Required
  }
  version = {
    db_version = "3.6.0" // Required
    ca_certificate = "" // If not set, uses default certificate from project
    ip_whitelist = "" // If not set, no whitelist is configured
  }
  configuration = {
    model = "oneshard"
    node_size_id = "a4"
    node_disk_size = 20
  }
}

// Example of a sharded deployment
resource "oasis_deployment" "my_sharded_deployment" {
  project = "190765139" // If set here, overrides project in provider
  location = {
    provider = "gcp" // Required
    region = "gcp-eu-west4" // Required
  }
  version = {
    db_version = "3.6.0" // Required
    ca_certificate = "" // If not set, uses default certificate from project
    ip_whitelist = "" // If not set, no whitelist is configured
  }
  configuration = {
    model = "sharded"
    node_size_id = "a4"
    node_disk_size = 20
    num_nodes = 5
  }
}
resource "oasis_deployment" "my_flexible_deployment" {
  project = "190765139" // If set here, overrides project in provider
  location = {
    provider = "gcp" // Required
    region = "gcp-eu-west4" // Required
  }
  version = {
    db_version = "3.6.0" // Required
    ca_certificate = "" // If not set, uses default certificate from project
    ip_whitelist = "" // If not set, no whitelist is configured
  }
  configuration = {
    model = "flexible"
    coordinator_memory_size = 3
    dbserver_memory_size = 8
    dbserver_disk_size = 64
    num_coordinators = 3
    num_dbservers = 5
  }
}
```

## Project Data Source

To define and use a project as data source, use the following terraform configuration:

```
# Find the latest available AMI that is tagged with Component = web
data "oasis_project" "my_project" {
  name = "MyProject"
  
}
```
