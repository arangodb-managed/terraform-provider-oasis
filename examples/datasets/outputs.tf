// Load in all the available datasets
data "oasis_example_datasets" "datasets" {
  // optionally define an organization id to list example datasets for which are
  // only available to that organiztaion. If you do not have access to said organization
  // this will just be ignored.
  organization = "_support" // can also be defined in provider section.
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_example_datasets.datasets
}