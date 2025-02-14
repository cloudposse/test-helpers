variable "cnt" {}

variable "namespace" {}

variable "stage" {}

variable "environment" {}

variable "tenant" {}

resource "null_resource" "test" {
  count = var.cnt
}
