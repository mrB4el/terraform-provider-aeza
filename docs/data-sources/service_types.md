---
page_title: "aeza_service_types Data Source - terraform-provider-aeza"
subcategory: ""
description: |-
  Get available Aeza service types.
---

# aeza_service_types (Data Source)

Use this data source to retrieve available service types.

## Example Usage

```terraform
data "aeza_service_types" "all" {}

output "service_types" {
  value = data.aeza_service_types.all.types
}
```

## Schema

### Read-Only

- `id` (String) The ID of this data source.
- `types` (List of Object) List of available service types.

The v2 API does not provide an `order` field for service types; its value is `0`.

