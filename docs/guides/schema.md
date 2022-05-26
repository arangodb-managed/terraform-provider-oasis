---
page_title: "Schema"
description: |-
    Guide to generating the schema for the Terraform Provider
---

# Schema

In order to see every configuration option for this plugin, either browse the code for the data source
you are interested in, or, once an Oasis Terraform configuration file is provided, take a look at the schema
with the following command:

```bash
terraform providers schema -json ./my_oasis_deployment | jq
```

Where `./my_oasis_deployment` is a folder which contains terraform configuration files.