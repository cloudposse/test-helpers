variable "region" {}

variable "attributes" {
  type    = list
  default = []
}

output "test" {
  value = "Hello, World"
}

output "test_list" {
  value = [ "a", "b", "c"]
}

output "test_map_of_objects" {
  value = {
    a = {
      b = "c"
    },
    d = {
      e = "f"
    }
  }
}

