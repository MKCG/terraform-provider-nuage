---
page_title: "project Resource - terraform-provider-nuage"
subcategory: ""
description: ""
---

# Resource `nuage_project`

## Example Usage

```terraform
resource "nuage_project" "project" {
  description = "projet infra"
  name = "infra00000000000"
}
```

## Argument Reference

- `description` - (Required) Description of the project
- `name` - (Required) Must be unique and match ^([a-z0-9]{16})$
