---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oasis_organization Data Source - terraform-provider-oasis"
subcategory: ""
description: |-
  
---

# oasis_organization (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `description` (String)
- `name` (String)

### Read-Only

- `created_at` (String)
- `id` (String) The ID of this resource.
- `is_deleted` (Boolean)
- `tier` (Set of Object) (see [below for nested schema](#nestedatt--tier))
- `url` (String)

<a id="nestedatt--tier"></a>
### Nested Schema for `tier`

Read-Only:

- `has_backup_uploads` (Boolean)
- `has_support_plans` (Boolean)
- `id` (String)
- `name` (String)
- `requires_terms_and_conditions` (Boolean)

