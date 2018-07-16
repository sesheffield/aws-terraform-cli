
# 
# Bucket for Terraform state files
#
resource "aws_s3_bucket" "tfstate_useast1" {
  bucket = "${var.project_name}-tfstate-useast1"
  acl    = "private"
  versioning {
    enabled = true
  }
  tags {
    Name    = "${var.project_name}-tfstate-useast1"
    Project = "${var.project_name}"
    Meta    = "author:${var.provisioner}"
  }
}
