# Example: Backup Data Source

This example shows how to use the Terraform Oasis provider to manage Backup Data Source in Oasis.

## Prerequisites

*This example uses syntax elements specific to Terraform version 0.13+ (tested on Terraform version 1.1.4).
It will not work out-of-the-box with Terraform 0.12.x and lower (deprecated by Terraform).*

## Environment variables
Please refer to [Main README](../../README.md) file for all the environment variables you might need.

## Example output:
```
backup = {
  "backup_policy_id" = "0syap73n4gqpbikn6oo5"
  "created_at" = "2022-05-26T08:00:10.835Z"
  "deployment_id" = "bszpx3zjnquchcwvil22"
  "description" = "This backup has been created by backup policy (scheduled backup): Default backup policy"
  "id" = "y2pcdn45pvrm8gld8zmv"
  "name" = "Created 2022-05-26T08:00:07Z"
  "url" = "/Organization/_support/Project/76253721/Deployment/bszpx3zjnquchcwvil22/Backup/y2pcdn45pvrm8gld8zmv"
}
```

## Instructions on how to run:
```
terraform init
terraform plan
terraform apply
```

To remove the resources created run:
```
terraform destroy
``` 