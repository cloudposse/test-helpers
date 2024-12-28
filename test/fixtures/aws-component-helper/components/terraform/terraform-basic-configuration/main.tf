variable "cnt" {}

variable "region" {}

resource "null_resource" "test" {
  count = var.cnt
}
