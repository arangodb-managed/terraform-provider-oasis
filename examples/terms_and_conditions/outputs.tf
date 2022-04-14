// Load in Oasis terms and conditions
data "oasis_terms_and_conditions" "test_terms_and_conditions" {
}

// Output the data after it has been synced.
output "datasets" {
  value = data.oasis_terms_and_conditions.test_terms_and_conditions
}