---
page_title: "aeza_recipes Data Source - terraform-provider-aeza"
subcategory: ""
description: |-
  Get available installation recipes.
---

# aeza_recipes (Data Source)

Use this data source to retrieve available installation recipes (pre-configured software stacks).

## Example Usage

```terraform
data "aeza_recipes" "all" {}

output "recipes" {
  value = data.aeza_recipes.all.recipes
}
```

## Schema

### Read-Only

- `id` (String) The ID of this data source.
- `recipes` (List of Object) List of available installation recipes.

The v2 API provides `operatingSystemSlug`, exposed as `operating_system`. The numeric `os` is resolved by matching that slug to the operating systems list; it is `0` if no match exists. The legacy `tags` and `targets` fields are empty when absent from the v2 response.

