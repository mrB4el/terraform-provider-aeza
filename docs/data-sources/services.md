---
page_title: "aeza_services Data Source - terraform-provider-aeza"
subcategory: ""
description: |-
  Get list of existing Aeza services.
---

# aeza_services (Data Source)

Use this data source to retrieve information about your existing Aeza services.

## Example Usage

```terraform
# Get all services
data "aeza_services" "all" {}

# Output services
output "my_services" {
  value = data.aeza_services.all.services
}

# Filter active services
output "active_services" {
  value = [
    for service in data.aeza_services.all.services :
    service if service.status == "active"
  ]
}
```

## Schema

### Read-Only

- `id` (String) The ID of this data source.
- `services` (List of Object) List of services with their configurations and status.
  Each service includes its primary `ip` in `ipv4_addresses` or `ipv6_addresses`. For services with the `ip` capability, the provider also fetches both network lists, adding two API requests per service.
  IPv6 subnets are exposed separately in `ipv6_subnets` and do not imply that individual addresses are assigned.
