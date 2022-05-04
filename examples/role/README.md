# Example: IAM Role

This example shows how to use the Terraform Oasis provider to create an IAM Role for a specific organization within Oasis.

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

__Note__: When creating a role please make sure to look in Oasis dashboard or use [oasisctl](https://github.com/arangodb-managed/oasisctl) to get the list of allowed permissions, if you specify a permission that is not allowed, you will not be able to create the role. Note that permissions is an optional field, it can be updated later.

To use `oasisctl` to get the list of permissions you can run:
```
oasisctl list permissions
```
This will output a list of permissions you can use for the `permissions` field when creating an `oasis_iam_role` resource.

To remove the resources created run:
```
terraform destroy
```