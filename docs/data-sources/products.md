---
page_title: "aeza_products Data Source - terraform-provider-aeza"
subcategory: ""
description: |-
  Get information about available Aeza products.
---

# aeza_products (Data Source)

Use this data source to get information about available Aeza products including configurations, prices, and service types.

## Example Usage

```terraform
# Get all products
data "aeza_products" "all" {}

# Output products count
output "products_count" {
  value = length(data.aeza_products.all.products)
}

# Filter VPS products
output "vps_products" {
  value = [
    for product in data.aeza_products.all.products :
    product if product.type == "vps"
  ]
}

# Group products by type
output "products_by_type" {
  value = {
    for product in data.aeza_products.all.products :
    product.type => product.name...
  }
}

# Find products with monthly pricing
output "monthly_priced_products" {
  value = [
    for product in data.aeza_products.all.products :
    product if can(product.prices["month"])
  ]
}
```

## Schema

### Read-Only

- `id` (String) The ID of this data source.
- `products` (List of Object) List of available products. Each product has the following attributes:
  - `id` (Number) Unique product identifier.
  - `name` (String) Product name (e.g., "MSK-1", "AMD Ryzen 9 5950X – 10 Gbps").
  - `type` (String) Product type (vps, dedicated, domain, storage, hicpu, etc.).
  - `group_id` (Number) Product group identifier.
  - `service_handler` (String) Service handler (vm6, manual, feru, s3, ispmgr, etc.).
  - `prices` (Map of Number) Prices for different billing periods (hour, month, year, half_year, quarter_year).

## Notes

### Product Types

Product types correspond to service types from [aeza_service_types](service_types.md) data source.

### Service Handlers

Service handlers determine the service management system:
- `vm6` - Cloud servers (VPS)
- `manual` - Dedicated servers
- `feru` - Domain registration
- `s3` - Object storage
- `ispmgr` - ISPManager licenses

### Pricing

Prices are the v2 API's `rawPrices`, specified in the smallest currency units (kopecks, cents, etc.).

