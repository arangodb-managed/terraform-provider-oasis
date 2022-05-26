# Example: Deployment

This example shows how to use the Terraform Oasis provider to create an Oasis Deployment. A deployment contains an ArangoDB, configured as you choose. You can have any number of deployments under one project.

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

You can lock a Deployment by specifying the lock as an option in the schema:
```terraform
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project = "" // Project id where deployment will be created
  name = "oasis_test_dep_tf"
  location {
    region = "gcp-europe-west4"
  }
  version {
    db_version = "3.9.1"
  }
  configuration {
    model = "oneshard"
  }
  notification_settings {
    email_addresses = [
      "test@arangodb.com"
    ]
  }
  locked = true
}
```
Note: if you run `terraform destroy` while the Deployment is locked, an error is shown, that's because you can't delete a locked Deployment.
To delete it you have to either remove the property or set `locked=false`:
```terraform
resource "oasis_deployment" "my_oneshard_deployment" {
  terms_and_conditions_accepted = "true"
  project = "" // Project id where deployment will be created
  name = "oasis_test_dep_tf"
  location {
    region = "gcp-europe-west4"
  }
  version {
    db_version = "3.9.1"
  }
  configuration {
    model = "oneshard"
  }
  notification_settings {
    email_addresses = [
      "test@arangodb.com"
    ]
  }
  locked = false
}
```
After running `terraform plan` and then `terraform apply --auto-approve` you update the Deployment to not be locked anymore. This way you can run `terraform destroy` without errors, deleting the Deployment.

To remove the resources created run:
```
terraform destroy
```