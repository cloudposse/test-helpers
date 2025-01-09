variable "region" {}

resource "null_resource" "test" {
  triggers = {
    time = timestamp()
  }
}
