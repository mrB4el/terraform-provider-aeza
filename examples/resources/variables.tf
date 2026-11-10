variable "aeza_api_key" {
  description = "API Key for Aeza API"
  type        = string
  sensitive   = true
  validation {
    condition     = length(var.aeza_api_key) > 0
    error_message = "API key must not be empty."
  }
}

variable "aeza_base_url" {
  description = "Base URL for Aeza API"
  type        = string
  default     = "https://my.aeza.net/api"
}

variable "product_name" {
  description = "Exact Aeza product name shown in the control panel"
  type        = string
}
