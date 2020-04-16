/*
 * This simple example will demostrate how to list all available example datasets.
 */
provider "oasis" {
  api_key_id     = ""
  api_key_secret = ""
}

data "oasis_example_datasets" "datasets" {
  // optionally define an organization id to list example datasets for which are
  // only available to that organiztaion. If you do not have access to said organization
  // this will just be ignored.
  organization_id = "" // can also be defined in provider section.
}