// Load in Oasis Cloud Providers
data "oasis_region" "test_oasis_region" {
  organization = "" // put your organization id here
  provider_id  = "" // put one of the cloud provider ids here (can be fetched from cloud_provider data source)
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_region.test_oasis_region
}