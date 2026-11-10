# examples/resources/main.tf
terraform {
  required_providers {
    aeza = {
      source  = "mrb4el/aeza"
      version = "0.1.0"
    }
  }
}

provider "aeza" {
  api_key  = var.aeza_api_key
  base_url = var.aeza_base_url
}

resource "aeza_service" "test_vps" {
  # неизменные параметры
  product_name = var.product_name
  os           = "ubuntu_2404"
  payment_term = "hour"
  wait_for_ip  = true
  # изменяемые параметры
  name         = "test-updated-vps"
  auto_prolong = false
}


resource "aeza_service_prolong" "extend" {
  service_id = aeza_service.test_vps.id

  # method = "balance" (по умолчанию)
  term       = "hour"
  term_count = 1
  # force  = false (по умолчанию)
}

resource "aeza_service_actions" "manage" {
  service_id = aeza_service.test_vps.id

  is_active = true # выключить сервис
  ## restart   = true   # перезагрузить (выполнится если is_active = true)
}
