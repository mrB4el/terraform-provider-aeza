# Terraform Provider for Aeza

**Status:** Work in Progress (v0.1.0)

Terraform provider for managing [Aeza](https://aeza.net/) hosting services through their API.

## Overview

This provider allows you to manage Aeza services (VPS, dedicated servers, etc.) using Terraform infrastructure as code.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for building from source)

## Installation

### From Terraform Registry

```hcl
terraform {
  required_providers {
    aeza = {
      source  = "mrb4el/aeza"
      version = "~> 0.1.0"
    }
  }
}

provider "aeza" {
  api_key  = var.aeza_api_key
  base_url = var.aeza_base_url  # Optional, defaults to https://my.aeza.net/api
}
```

### Building from Source

```bash
git clone https://github.com/mrB4el/terraform-provider-aeza
cd terraform-provider-aeza
make build
```

## Configuration

The provider requires an API key for authentication:

```hcl
provider "aeza" {
  api_key = "your-api-key-here"
}
```

Or set via environment variable:

```bash
export AEZA_API_KEY="your-api-key-here"
export AEZA_BASE_URL="https://my.aeza.net/api"  # Optional
```

### Getting an API Key

1. Log in to your [Aeza control panel](https://my.aeza.net/)
2. Navigate to the API section
3. Generate a new API token

## Available Resources

### Data Sources

- `aeza_products` - List available products
- `aeza_services` - List active services
- `aeza_os_list` - List available operating systems
- `aeza_service_groups` - List service groups
- `aeza_service_types` - List available service types
- `aeza_recipes` - List available installation recipes

### Resources

- `aeza_service` - Create and manage services
- `aeza_service_prolong` - Extend service subscription
- `aeza_service_actions` - Manage service actions (start, stop, reboot)

## Example Usage

See the [examples](./examples) directory for complete usage examples:

- [basic](./examples/basic) - Basic provider usage
- [data-sources](./examples/data-sources) - Using data sources
- [resources](./examples/resources) - Creating and managing resources

### Quick Example

```hcl
data "aeza_products" "available" {}

resource "aeza_service" "my_vps" {
  product_id = 123
  os_id      = 456
  period     = "month"
  count      = 1
}
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Installing Locally

```bash
make install
```

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.

## Organization

Maintained by [SCINFRA](https://github.com/scinfra-pro)

## Acknowledgments

- Thanks to [SCINFRA](https://github.com/scinfra-pro) for the original development.
- Thanks to [supercakecrumb](https://github.com/supercakecrumb) for major improvements.

