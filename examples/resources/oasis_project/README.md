# Example: Project

This example shows how to use the Terraform Oasis provider to create an Oasis Project within an organization.

## Prerequisites

*This example uses syntax elements specific to Terraform version 0.13+ (tested on Terraform version 1.1.4).
It will not work out-of-the-box with Terraform 0.12.x and lower (deprecated by Terraform).*

## Environment variables
Please refer to [Main README](../../README.md) file for all the environment variables you might need.

## Instructions on how to run:
```
terraform init
terraform plan
terraform apply
```

You can lock a project by specifying the lock as an option in the schema:
```terraform
resource "oasis_project" "oasis_test_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
  locked      = true
}
```
Note: if you run `terraform destroy` while the project is locked, an error is shown, that's because you can't delete a locked project.
To delete it you have to either remove the property or set `lock=false`:
```terraform
resource "oasis_project" "oasis_test_project" {
  name        = "Terraform Oasis Project"
  description = "A test Oasis project within an organization from the Terraform Provider"
  locked      = false
}
```
After running `terraform plan` and then `terraform apply --auto-approve` you update the project to not be locked anymore. This way you can run `terraform destroy` without errors, this way deleting the project.

To remove the resources created run:
```
terraform destroy
``` 