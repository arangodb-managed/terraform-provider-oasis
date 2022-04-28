# Example: Backup Policy

This example shows how to use the Terraform Oasis provider to create a backup policy for an Oasis Deployment.

## Prerequsites

*This example uses syntax elements specific to Terraform version 0.13+ (tested on Terraform version 1.1.4).
It will not work out-of-the-box with Terraform 0.12.x and lower.*

## Environment variables
Please refer to [Main README](../../README.md) file for all the environment variables you might need.

## Instructions on how to run:
```
terraform init
terraform plan
terraform apply
```

You can lock a Backup Policy by specifying the lock as an option in the schema:
```terraform
resource "oasis_backup_policy" "my_backup_policy" {
  name               = "Test Backup Policy"
  description        = "Test Description"
  email_notification = "FailureOnly"
  deployment_id      = "" // Provide Deployment ID here.
  retention_period_hour   = 120
  upload             = true
  schedule {
    type = "Monthly"
    monthly {
      day_of_month = 12
      schedule_at {
        hours    = 15
        minutes  = 10
        timezone = "UTC"
      }
    }
  }
  locked = true
}
```
Note: if you run `terraform destroy` while the Backup Policy is locked, an error is shown, that's because you can't delete a locked Backup Policy.
To delete it you have to either remove the property or set `lock=false`:
```terraform
resource "oasis_backup_policy" "my_backup_policy" {
  name               = "Test Backup Policy"
  description        = "Test Description"
  email_notification = "FailureOnly"
  deployment_id      = "" // Provide Deployment ID here.
  retention_period_hour   = 120
  upload             = true
  schedule {
    type = "Monthly"
    monthly {
      day_of_month = 12
      schedule_at {
        hours    = 15
        minutes  = 10
        timezone = "UTC"
      }
    }
  }
  locked = false
}
```
After running `terraform plan` and then `terraform apply --auto-approve` you update the Backup Policy to not be locked anymore. This way you can run `terraform destroy` without errors, deleting the Backup Policy.

To remove the resources created run:
```
terraform destroy
```