variable "default_tags" {
  type    = map(string)
  default = {}
}

variable "region" {
}

variable "revision" {
  type = string
}

locals {
  revision = join("-", concat(module.this.attributes, [var.revision]))
}

resource "terraform_data" "replacement" {
  input = local.revision
}

output "revision" {
  value = local.revision
}
