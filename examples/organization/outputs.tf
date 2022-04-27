// Load in an Oasis organization
data "oasis_organization" "oasis_test_organization" {
  id = oasis_organization.oasis_test_organization.id // provide your organization id here
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_organization.oasis_test_organization
}