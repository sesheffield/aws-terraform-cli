#
# Core state
#
terraform {
	backend "s3" {
	bucket = "terraform-tfstate-us-east-1"
	region = "us-east-1"
	key = "terraform-apps-ops.tfstate"
	acl = "private"
  }
}
