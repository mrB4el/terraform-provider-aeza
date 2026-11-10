---
page_title: "aeza_os_list Data Source - terraform-provider-aeza"
subcategory: ""
description: |-
  Get available operating systems.
---

# aeza_os_list (Data Source)

Use this data source to retrieve available operating systems for services.

## Example Usage

```terraform
data "aeza_os_list" "all" {}

output "operating_systems" {
  value = data.aeza_os_list.all.operating_systems
}
```

## Schema

### Read-Only

- `id` (String) The ID of this data source.
- `operating_systems` (List of Object) List of available operating systems.

The legacy `repository` and `targets` fields are empty when absent from the v2 response.

