# examples/basic/main.tf
terraform {
  required_providers {
    aeza = {
      source  = "mrb4el/aeza"
      version = "0.1.0"
    }
  }
}

provider "aeza" {
  api_key = var.aeza_api_key
  base_url = var.aeza_base_url
}


# Data source для типов услуг
data "aeza_service_types" "all" {}

# Data source для получения списка продуктов
data "aeza_products" "all" {}

# Data source для услуг пользователя
data "aeza_services" "all" {}

# Data source для групп услуг, есть параметр type
data "aeza_service_groups" "all" {
  # Optional: filter by specific service type
  # service_type = "vps"
}

# Data source для операционных систем
data "aeza_os_list" "available" {}

# Data source для рецептов установки
data "aeza_recipes" "all" {}
