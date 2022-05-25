// Load in all the available datasets
data "oasis_example_datasets" "datasets" {
  // optionally define an organization id to list example datasets for which are
  // only available to that organization. If you do not have access to said organization
  // this will just be ignored.
  organization = "" // organization id
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_example_datasets.datasets
}