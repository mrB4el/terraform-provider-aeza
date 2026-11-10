---
page_title: "Provider: Aeza"
description: |-
  The Aeza provider is used to interact with Aeza cloud resources.
---

# Aeza Provider

The Aeza provider allows Terraform to manage [Aeza](https://aeza.net/) cloud resources including VPS, dedicated servers, and other services.

## Example Usage

```terraform
terraform {
  required_providers {
    aeza = {
      source  = "mrb4el/aeza"
      version = "~> 0.1.0"
    }
  }
}

provider "aeza" {
  api_key = var.aeza_api_key
  # base_url = "https://my.aeza.net/api"  # Optional: custom API endpoint
}
```

## Authentication

The Aeza provider requires an API key for authentication. You can obtain your API key from the [Aeza control panel](https://my.aeza.net/).

### API Key

The API key can be provided in several ways:

1. **Provider configuration** (recommended for testing):
   ```terraform
   provider "aeza" {
     api_key = "your-api-key-here"
   }
   ```

2. **Environment variable** (recommended for production):
   ```bash
   export AEZA_API_KEY="your-api-key-here"
   ```

3. **Variable**:
   ```terraform
   variable "aeza_api_key" {
     type      = string
     sensitive = true
   }
   
   provider "aeza" {
     api_key = var.aeza_api_key
   }
   ```

## Schema

### Optional

- `api_key` (String, Sensitive) The API key for authenticating with Aeza API. Can also be sourced from `AEZA_API_KEY` environment variable.
- `base_url` (String) Base URL for the Aeza API. Defaults to `https://my.aeza.net/api`.

## Resources

- [aeza_service](resources/service.md) - Manage Aeza services (VPS, dedicated servers, etc.)
- [aeza_service_actions](resources/service_actions.md) - Perform actions on services (start, stop, reboot, etc.)
- [aeza_service_prolong](resources/service_prolong.md) - Extend service subscription period

## Data Sources

- [aeza_products](data-sources/products.md) - Get information about available products
- [aeza_services](data-sources/services.md) - Get list of existing services
- [aeza_service_types](data-sources/service_types.md) - Get available service types
- [aeza_service_groups](data-sources/service_groups.md) - Get service groups
- [aeza_os_list](data-sources/os_list.md) - Get available operating systems
- [aeza_recipes](data-sources/recipes.md) - Get installation recipes

