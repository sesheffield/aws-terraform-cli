
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
  default = "1234563562"
}

# AWS REGION
variable "aws_region" {
  default = "us-east-1"
}

# AWS PROJECT NAME 
variable "project_name" {
  default = "ingersollrand"
}

# AWS DEV ENVIRONMENT 
variable "dev_environment" {
  default = "development"
}

# AWS DEV ENV
variable "dev_env" {
  default = "dev"
}

# AWS STG ENVIRONMENT 
variable "stg_environment" {
  default = "staging"
}

# AWS STG ENV
variable "stg_env" {
  default = "stg"
}

# AWS PRD ENVIRONMENT 
variable "prd_environment" {
  default = "production"
}

# AWS PRD ENV
variable "prd_env" {
  default = "prd"
}

# AWS GRD ENV
variable "grd_env" {
  default = "grd"
}

# AWS PROVISIONER
variable "provisioner" {
  default = "terraform"
}
