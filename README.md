# ArangoDB Oasis Terraform Provider

![ArangoDB Oasis](https://cloud.arangodb.com/assets/logos/arangodb-oasis-logo-whitebg-right.png)

## Project status: Preview

Note that this provider is currently in preview.

**That means that its API may still change.**

We welcome your feedback!

## Maintainers

This provider plugin is maintained by the team at [ArangoDB](https://www.arangodb.com/).

## Installation

Downloading the [latest released binaries](https://github.com/arangodb-managed/terraform-provider-oasis/releases),
extract the zip archive and install the binary for your platform in your preferred location.

Or to build from source, run:

```bash
git clone https://github.com/arangodb-managed/terraform-provider-oasis.git
make
```

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+ (recommend 1.1.4+) 
- [Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

## Usage

### Retrieving API Keys

Api keys can be generated and viewed under the user's dashboard view on the API Keys tab.
On a logged in view, navigate to [API Keys](https://cloud.arangodb.com/dashboard/user/api-keys) and hit the button
labeled `New API key`. This will generate a set of keys which can be used with ArangoDB's public API.

### Retrieving relevant IDs

After a key has been obtained, the relevant API can be called to list organizations, projects etc.

## Running Acceptance Tests

In order to run acceptance tests, the following make target needs to be executed:

```bash
make test-acc
```

It is recommended that on a schema addition / deprecation and general larger refactorings the acceptance tests are
executed. *NOTE* that these tests create real deployments, projects and organizations.

All of them require the following two environment properties to be set:

```bash
export OASIS_API_KEY_ID=<your_key_id>
export OASIS_API_KEY_SECRET=<your_key_secret>
```

In addition, those properties might be needed:
```bash
export OASIS_ENDPOINT=<oasis_endpoint>, 
export OASIS_TEST_ORGANIZATION_ID=<organization_id>, 
export OASIS_PROJECT=<oasis_project>
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

