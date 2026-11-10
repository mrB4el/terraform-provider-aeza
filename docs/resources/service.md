---
page_title: "aeza_service Resource - terraform-provider-aeza"
subcategory: ""
description: |-
  Manages an Aeza service (VPS, dedicated server, etc.).
---

# aeza_service (Resource)

Manages an Aeza service including VPS, Hi-CPU servers, and dedicated servers.

The recommended selector is `product_name`, using an exact name from the generated [Aeza product-name catalog](../product-catalog.md). The numeric `product_id` remains available as a fallback when the API returns duplicate names.

You can also inspect names directly from Terraform before creating anything:

```terraform
data "aeza_products" "available" {}

output "available_vps_product_names" {
  value = sort(distinct([
    for product in data.aeza_products.available.products : product.name
    if product.type == "vps"
  ]))
}
```

## Example Usage

### Basic VPS

```terraform
variable "product_name" {
  description = "Exact product name shown in the Aeza control panel"
  type        = string
}

resource "aeza_service" "web_server" {
  product_name = var.product_name
  name         = "web-server"
  os           = "ubuntu_2404"
  payment_term = "month"
  wait_for_ip  = true
}

output "web_server_addresses" {
  value = {
    ipv4 = tolist(aeza_service.web_server.ipv4_addresses)
    ipv6 = tolist(aeza_service.web_server.ipv6_addresses)
  }
}
```

### VPS with Custom Configuration

```terraform
resource "aeza_service" "app_server" {
  # Numeric IDs remain supported when a product name is not unique.
  product_id   = 181
  name         = "app-server"
  os           = "ubuntu_2404"
  recipe       = "docker"
  auto_prolong = true
}
```

## Schema

### Required

- `name` (String) Name of the service.
- One of:
  - `product_name` (String) Exact product name. Matching is case-insensitive. If several products share the name, use `product_id`.
  - `product_id` (Number) Numeric product ID.

### Optional

- `payment_term` (String) Billing period. Defaults to `month`.
- `auto_prolong` (Boolean) Enable automatic renewal. Defaults to `true`.
- `os` (String) Operating-system slug, such as `ubuntu_2404`.
- `recipe` (String) Installation recipe slug.
- `iso_url` (String) Custom installation ISO URL.
- `wait_for_ip` (Boolean) Wait up to the resource create timeout for at least one IPv4 or IPv6 address. Defaults to `true`.

### Read-Only

- `id` (String) The ID of the service.
- `status` (String) Current status of the service.
- `ip` (String) Primary IP address, kept for backward compatibility.
- `ipv4_addresses` (Set of String) All IPv4 addresses assigned to the service.
- `ipv6_addresses` (Set of String) All IPv6 addresses assigned to the service.
- `created_at` (String) Service creation timestamp.

`product_name` and the resolved `product_id` are both retained in state, so downstream modules can use either value.

## Keeping product names current

Product names are owned by Aeza and are not hardcoded in the provider. The repository includes `cmd/product-catalog`, which reads the current catalog from the Aeza API and deterministically updates `docs/product-catalog.md`:

```bash
AEZA_API_KEY="$AEZA_API_KEY" go run ./cmd/product-catalog \
  --output docs/product-catalog.md
```

The scheduled **Update Aeza product catalog** GitHub Actions workflow runs this command daily. It opens or updates a pull request from `github-actions[bot]` only when the generated catalog changes. Configure the repository secret `AEZA_API_KEY`; `AEZA_BASE_URL` is optional.

The repository must also allow GitHub Actions to create pull requests under **Settings → Actions → General → Workflow permissions**. The workflow never prints the API response or API key to the Actions log.

## Import

Services can be imported using the service ID:

```bash
terraform import aeza_service.example 12345
```
