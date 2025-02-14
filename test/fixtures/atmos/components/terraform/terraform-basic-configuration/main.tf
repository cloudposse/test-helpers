variable "cnt" {}

variable "namespace" {}

variable "stage" {}

resource "null_resource" "test" {
  count = var.cnt
}
