# Example: Current User Data Source

This example shows how to use the Terraform Oasis provider to manage Current User Data Source in Oasis.

## Prerequisites

*This example uses syntax elements specific to Terraform version 0.13+ (tested on Terraform version 1.1.4).
It will not work out-of-the-box with Terraform 0.12.x and lower (deprecated by Terraform).*

## Environment variables
Please refer to [Main README](../../README.md) file for all the environment variables you might need.

## Example output
```
datasets = {
  "email" = "test@arangodb.com"
  "id" = "google-oauth2|781471512049238536110"
  "name" = "Test Arango"
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