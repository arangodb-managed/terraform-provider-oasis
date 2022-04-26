// Load in an Oasis backup of a deployment
data "oasis_backup" "oasis_test_backup" {
  id = oasis_backup.my_backup.id // Backup ID
}

// Output the data after it has been synced.
output "backup" {
  value = data.oasis_backup.oasis_test_backup
}