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
  count = module.this.enabled ? 1 : 0
  input = local.revision
}

output "revision" {
  value = module.this.enabled ? local.revision : null
}
