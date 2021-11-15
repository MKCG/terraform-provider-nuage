variable "nuage_auth_organization" {
  type = string
  sensitive = true
  description = "Nuage organization"
}

variable "nuage_auth_name" {
  type = string
  sensitive = true
  description = "Nuage username"
}

variable "nuage_auth_password" {
  type = string
  sensitive = true
  description = "Nuage password"
}
