variable "cnt" {
  type        = number
  description = "Number of null_resource instances to create"
  validation {
    condition     = var.cnt > 0
    error_message = "Count must be greater than 0"
  }
}

variable "region" {
  type        = string
  description = "AWS region for the resources"
  validation {
    condition     = can(regex("^[a-z]{2}(-[a-z]+)?-[1-2]$", var.region))
    error_message = "Region must be a valid AWS region name (e.g., us-west-2)"
  }
}

variable "attributes" {
  type        = list(string)
  default     = []
  description = "Additional attributes to add to the resources"
}

resource "null_resource" "test" {
  count = var.cnt
}
