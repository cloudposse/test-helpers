variable "seed" {
  type = string
}

resource "random_pet" "test" {
}

locals {
  pet_plus_seed = "${var.seed}-${random_pet.test.id}"
}

output "name" {
  value = local.pet_plus_seed
}
