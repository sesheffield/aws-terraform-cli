
# AWS PROVIDER
provider "aws" {
  region = "${var.aws_region}"
  allowed_account_ids = [
    "${var.aws_account_id}",
  ]
  version = "~> 1.3"
}

# AWS ACCOUNT ID
variable "aws_account_id" {
  default = "12344346456"
}

# AWS REGION
variable "aws_region" {
  default = "us-east-1"
}

# AWS PROJECT NAME 
variable "project_name" {
  default = "terraform"
}

# AWS PROVISIONER
variable "provisioner" {
  default = "terraform"
}
