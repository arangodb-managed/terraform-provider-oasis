# Example: Oasis Cloud Providers

This example shows how to use the Terraform Oasis provider to get the list of support cloud providers in Oasis.

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

Example output: 
```hcl
datasets = {
  "organization" = "_support"
  "providers" = tolist([
    {
      "id" = "aks"
      "name" = "Microsoft Azure"
    },
    {
      "id" = "aws"
      "name" = "Amazon Web Services"
    },
    {
      "id" = "gcp"
      "name" = "Google Compute Platform"
    },
  ])
}
```

To remove the resources created run:
```
terraform destroy
``` 