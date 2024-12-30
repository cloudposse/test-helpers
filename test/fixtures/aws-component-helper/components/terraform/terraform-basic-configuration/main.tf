variable "cnt" {}

variable "region" {}

variable "attributes" {
  type    = list
  default = []
}

resource "null_resource" "test" {
  count = var.cnt
}
