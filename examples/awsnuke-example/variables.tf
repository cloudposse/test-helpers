variable "region" {
  type        = string
  description = "AWS region"
}

variable "default_tags" {
  type        = map(string)
  description = "A map of tags to add to every resource"
  default     = {}
}
