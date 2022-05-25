# Example: Organization

This example shows how to use the Terraform Oasis provider to get data about an Oasis Organization. An organization typically represents a (commercial) entity such as a company, a company division, an institution or a non-profit organization.

## Prerequsites

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

You can lock an organization by specifying it as an option in the schema:
```terraform
resource "oasis_organization" "oasis_test_organization" {
  name        = "Terraform Oasis Organization"
  description = "A test Oasis organization from Terraform Provider"
  lock = true 
}
```
Note: if you run `terraform destroy` while the organization is locked an error is shown, that's because you can't delete a locked organization.
To delete it you have to either remove the property or set `lock=false`:
```terraform
resource "oasis_organization" "oasis_test_organization" {
  name        = "Terraform Oasis Organization"
  description = "A test Oasis organization from Terraform Provider"
  lock = false 
}
```
After running `terraform plan` and then `terraform apply --auto-approve` you update the organization to not be locked anymore. This way you can run `terraform destroy` without errors.

To remove the resources created run:
```
terraform destroy
``` 