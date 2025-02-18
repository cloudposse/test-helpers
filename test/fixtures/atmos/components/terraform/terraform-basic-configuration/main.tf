variable "cnt" {}

variable "namespace" {
  type        = string
  description = "ID element. Usually an abbreviation of the organization name, e.g., 'cp' for Cloud Posse"
}

variable "stage" {
  type        = string
  description = "ID element. Usually used to indicate role, e.g., 'prod', 'staging', 'source', 'build', 'test', 'deploy', 'release'"
}

variable "environment" {
  type        = string
  description = "ID element. Usually used to indicate type of environment, e.g., 'uw2', 'us-west-2', OR 'prod', 'staging', 'dev', 'UAT'"
}

variable "tenant" {
  type        = string
  description = "ID element. Used to indicate tenant or group, e.g., 'cloud', 'internal', 'restricted', 'shared'"
}
resource "null_resource" "test" {
  count = var.cnt
}
