terraform {
  required_providers {
    nuage = {
      source  = "nua.ge/terraform/nuage"
    }
  }
  required_version = "~> 1.0.3"
}

provider "nuage" {
  name = var.nuage_auth_name
  password = var.nuage_auth_password
  organization = var.nuage_auth_organization
}

data "nuage_flavor" "web_flavor" {
  core = 1
  ram  = 2
  disk = 100
}

data "nuage_flavor" "db_flavor" {
  core = 2
  ram  = 2
  disk = 100
}

data "nuage_image" "web_image" {
  os_name = "ubuntu"
  os_version = "21.04 (Hirsute Hippo)"
}

data "nuage_image" "db_image" {
  os_name = "ubuntu"
  os_version = "20.04 LTS (Focal Fossa)"
}

resource "nuage_keypair" "keypair" {
  description = "ssh key"
  public_key = file("~/.ssh/id_rsa.pub")
  is_default = true
}

resource "nuage_project" "project" {
  description = "projet infra"
  name = "infra00000000000"
}

resource "nuage_project" "project-bis" {
  description = "projet infra"
  name = "infra00000000001"
}

resource "nuage_server" "prod-web-1" {
  description = "web server"
  name    = "prod-web-1"
  project = nuage_project.project.id
  flavor  = data.nuage_flavor.web_flavor.id
  image   = data.nuage_image.web_image.id
  keypair = nuage_keypair.keypair.id
}

resource "nuage_server" "prod-web-2" {
  description = "web server"
  name    = "prod-web-2"
  project = nuage_project.project.id
  flavor  = data.nuage_flavor.web_flavor.id
  image   = data.nuage_image.web_image.id
  keypair = nuage_keypair.keypair.id
}

resource "nuage_server" "prod-db-1" {
  description = "web server"
  name    = "prod-db-1"
  project = nuage_project.project.id
  flavor  = data.nuage_flavor.db_flavor.id
  image   = data.nuage_image.db_image.id
  keypair = nuage_keypair.keypair.id
}

# resource "nuage_security_rule" "sec-rule" {
#   direction = "in"
#   protocol = "tcp"
#   port_min = 80
#   port_max = 80
#   remote = "192.168.0.1/32"
# }

# resource "nuage_security_rule" "sec-rule-bis" {
#   direction = "in"
#   protocol = "tcp"
#   port_min = 443
#   port_max = 443
#   remote = "192.168.0.1/32"
# }

# resource "nuage_security_group" "sec-group" {
#   name = "internal"
#   description = "Used by internal users"
#   rules = [
#     nuage_security_rule.sec-rule.id,
#     nuage_security_rule.sec-rule-bis.id
#   ]
# }

# output "sg" {
#   value = nuage_security_group.sec-group
# }
