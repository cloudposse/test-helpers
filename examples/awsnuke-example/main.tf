provider "aws" {
  region = var.region
  default_tags { tags = var.default_tags }
}

resource "aws_s3_bucket" "test" {
  bucket = module.this.id
}
