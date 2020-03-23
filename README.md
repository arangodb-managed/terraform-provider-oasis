# ArangoDB Oasis Terraform Provider

![ArangoDB Oasis](https://cloud.arangodb.com/assets/logos/arangodb-oasis-logo-whitebg-right.png)

**Note that this provider is currently in preview**.

That means that its API may still change.

We welcome your feedback!

## Maintainers

This provider plugin is maintained by the team at [ArangoDB](https://www.arangodb.com/).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.10.x
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Usage

### Retrieving API Keys

Api keys can be generated and viewed under the user's dashboard view on the API Keys tab.
On a logged in view, navigate to [API Keys](https://cloud.arangodb.com/dashboard/user/api-keys) and hit the button
labeled `New API key`. This will generate a set of keys which can be used with ArangoDB's public API.

### Retrieving relevant IDs

After a key has been obtained, the relevant API can be called to list organizations, projects etc.

## Configuration

The following is an example of a terraform deployment configuration:

```hcl
provider "oasis" {
  api_key_id = "xx"
  api_key_secret  = "xxx"
  organization = "190765105"
  project = "foo"
}

// Example of oneshard deployment
resource "oasis_deployment" "my_oneshard_deployment" {
  project = "190765139" // If set here, overrides project in provider
  location {
    region = "gcp-europe-west4" // Required
  }
  version {
    db_version = "3.6.0" // Required
  }

  security { // this section is optional
    ca_certificate = "" // If not set, uses default certificate from project
    ip_whitelist = "" // If not set, no whitelist is configured
  }

  configuration {
    model = "oneshard"
    node_size_id = "a4"
    node_disk_size = 20
  }
}

// Example of oneshard deployment without node specification
resource "oasis_deployment" "my_oneshard_deployment" {
  project = "190765139" // If set here, overrides project in provider
  location {
    region = "gcp-europe-west4" // Required
  }

  version {
    db_version = "3.6.0" // Required
  }

  security { // this section is optional
    ca_certificate = "" // If not set, uses default certificate from project
    ip_whitelist = "" // If not set, no whitelist is configured
  }

  configuration {
    model = "oneshard" // the smallest will be selected that's available in the given region
  }
}

// Example of a sharded deployment
resource "oasis_deployment" "my_sharded_deployment" {
  project = "190765139" // If set here, overrides project in provider
  location {
    region = "gcp-eu-west4" // Required
  }

  version {
    db_version = "3.6.0" // Required
  }

  security { // this section is optional
    ca_certificate = "" // If not set, uses default certificate from project
    ip_whitelist = "" // If not set, no whitelist is configured
  }

  configuration {
    model = "sharded"
    node_size_id = "a4"
    node_disk_size = 20
    num_nodes = 5
  }
}
```

## Data sources

### Project Data Source

To define and use a project as data source, consider the following terraform configuration:

```hcl
data "oasis_project" "my_project" {
  name = "MyProject"
  id = "123456789"
}

resource "oasis_deployment" "my_flexible_deployment" {
  project = data.oasis_project.my_project.id
}
```

### Organization Data Source

To define and use an organization as data source, consider the following terraform configuration:

```hcl
data "oasis_organization" "my_organization" {
  name = "MyOrganization"
  id = "123456789"
}

resource "oasis_deployment" "my_flexible_deployment" {
  organization = data.oasis_organization.my_organization.id
  ...
}
```

## Running Acceptance Tests

In order to run acceptance tests, the following make target needs to be executed:

```bash
make test-acc
```

It is recommended that on a schema addition / deprecation and general larger refactorings the acceptance tests are
executed. *NOTE* that these tests create real deployments, projects and organizations.

Some of them may require additional environment properties to work. I.e.:

```dotenv
OASIS_TEST_ORGANIZATION_ID=123456789
```

All of them require the following two environment properties to be set:

```dotenv
OASIS_API_KEY_ID=<your_key_id>
OASIS_API_KEY_SECRET=<your_key_secret>
```

## Examples

For further examples, please take a look under [Examples](./examples) folder.

## Schema

In order to see every configuration option for this plugin, either browse the code for the data source
you are interested in, or, once an oasis terraform configuration file is provided, take a look at the schema
with the following command:

```bash
terraform providers schema -json ./my_oasis_deployment | jq
```

Where `./my_oasis_deployment` is a folder which contains terraform configuration files.

## Links

- Website: https://cloud.arangodb.com/
- Slack: https://slack.arangodb.com/

