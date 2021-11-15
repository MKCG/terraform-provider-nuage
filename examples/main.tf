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

# resource "nuage_keypair" "keypair" {
#   description = "ssh key"
#   public_key = "abc"
#   user = "mkcg"  
# }

resource "nuage_project" "project" {
  description = "projet infra"
  name = "infra00000000000"
}

resource "nuage_project" "project-bis" {
  description = "projet infra"
  name = "infra00000000001"
}

resource "nuage_project" "project-ter" {
  description = "projet infra"
  name = "infra00000000002"
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
