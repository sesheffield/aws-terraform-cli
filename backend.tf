#
# Core state
#
terraform {
  backend "s3" {
	bucket = "ingersollrand-tfstate-useast1"
	region = "us-east-1"
	key = "ingersollrand-apps-ops.tfstate"
	encrypt = "true"
	acl = "private"
  }
}
