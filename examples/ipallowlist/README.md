# Example: IPAllowlist

This example shows how to use the Terraform Oasis provider to create an IPAllowlist for a specific project within Oasis.

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

You can lock an IP Allow List by specifying the lock as an option in the schema:
```terraform
resource "oasis_ipallowlist" "my_iplist" {
  name        = "Terraform IP Allowlist"
  description = "IP Allowlist description"
  cidr_ranges = ["1.2.3.4/32", "111.11.0.0/16"]
  project     = oasis_project.my_project.id
  locked      = true
}
```
Note: if you run `terraform destroy` while the IP Allow List is locked, an error is shown, that's because you can't delete a locked IP Allow List.
To delete it you have to either remove the property or set `lock=false`:
```terraform
resource "oasis_ipallowlist" "my_iplist" {
  name        = "Terraform IP Allowlist"
  description = "IP Allowlist description"
  cidr_ranges = ["1.2.3.4/32", "111.11.0.0/16"]
  project     = oasis_project.my_project.id
  locked      = false
}
```
After running `terraform plan` and then `terraform apply --auto-approve` you update the IP Allow List to not be locked anymore. This way you can run `terraform destroy` without errors, deleting the IP Allow List.

To remove the resources created run:
```
terraform destroy
```