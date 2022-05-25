// Load in an Oasis project within an organization
data "oasis_project" "oasis_test_project" {
  id = "" // Provide your project ID here
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_project.oasis_test_project
}