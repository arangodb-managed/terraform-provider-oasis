# Example: Cloud Provider Region Data Source

This example shows how to use the Terraform Oasis provider to manage Region Data Source in Oasis.

## Prerequisites

*This example uses syntax elements specific to Terraform version 0.13+ (tested on Terraform version 1.1.4).
It will not work out-of-the-box with Terraform 0.12.x and lower (deprecated by Terraform).*

## Environment variables
Please refer to [Main README](../../README.md) file for all the environment variables you might need.

## Example output
```
datasets = {
  "id" = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  "organization" = "_support"
  "provider_id" = "aks"
  "regions" = tolist([
    {
      "available" = true
      "id" = "aks-canadacentral"
      "location" = "Central Canada, Toronto"
      "provider_id" = "aks"
    },
    {
      "available" = true
      "id" = "aks-eastus2"
      "location" = "East US, Virginia"
      "provider_id" = "aks"
    },
    {
      "available" = true
      "id" = "aks-japaneast"
      "location" = "Japan East"
      "provider_id" = "aks"
    },
    {
      "available" = true
      "id" = "aks-southeastasia"
      "location" = "Southeast Asia, Singapore"
      "provider_id" = "aks"
    },
    {
      "available" = true
      "id" = "aks-uksouth"
      "location" = "UK, London"
      "provider_id" = "aks"
    },
    {
      "available" = true
      "id" = "aks-westeurope"
      "location" = "West Europe, Netherlands"
      "provider_id" = "aks"
    },
    {
      "available" = true
      "id" = "aks-westus2"
      "location" = "West US, Washington"
      "provider_id" = "aks"
    },
  ])
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