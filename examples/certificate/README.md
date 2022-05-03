# Example: Certificate

This example shows how to use the Terraform Oasis provider to create a Certificate for a specific project within Oasis.

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

You can lock a Certificate by specifying the lock as an option in the schema:
```terraform
resource "oasis_certificate" "my_oasis_cert" {
  name        = "Terraform certificate"
  description = "Certificate description."
  project     = oasis_project.oasis_test_project.id
  locked      = true
}
```
Note: if you run `terraform destroy` while the Certificate is locked, an error is shown, that's because you can't delete a locked CA Certificate.
To delete it you have to either remove the property or set `locked=false`:
```terraform
resource "oasis_certificate" "my_oasis_cert" {
  name        = "Terraform certificate"
  description = "Certificate description."
  project     = oasis_project.oasis_test_project.id
  locked      = false
}
```

To remove the resources created run:
```
terraform destroy
```