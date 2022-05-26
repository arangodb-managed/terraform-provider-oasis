---
page_title: "Getting started with Oasis Terraform Provider"
description: |-
    Guide to getting started with the ArangoDB Oasis Terraform Provider
---

# Getting started with Oasis Terraform Provider


### Retrieving API Keys

You will need API Keys to interact with the Terraform Provider. Api keys can be generated and viewed under the user's dashboard view on the API Keys tab.
On a logged in view, navigate to [API Keys](https://cloud.arangodb.com/dashboard/user/api-keys) and hit the button
labeled `New API key`. This will generate a set of keys which can be used with ArangoDB's public API.

### Retrieving relevant IDs

After a key has been obtained, the relevant API can be called to manage different resources.

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
