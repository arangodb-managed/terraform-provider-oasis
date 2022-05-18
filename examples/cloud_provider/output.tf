// Load in Oasis Cloud Providers
data "oasis_cloud_provider" "test_oasis_cloud_providers" {
  organization = "" // put your organization id here
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_cloud_provider.test_oasis_cloud_providers
}