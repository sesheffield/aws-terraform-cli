package cmd

// Vars ...
type Vars struct {
	Required
	Optional
}

// Required ...
type Required struct {
	StateBucket  string
	Region       string
	ConcatRegion string
	AccountID    int64
	ProjectName  string
	Environment  string
	S3BucketName string
}

// Optional ...
type Optional struct {
	ACL                  string
	Versioning           string
	Encrypt              string
	ServerSideEncryption string
}

const aclTF = `acl = "private"`

const versioningTF = `
	versioning {
		enabled = true
	}
`

const encryptTF = `encrypt = "true"`

const serverSideEncryptionTF = `
	server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
      }
    }
  }
`
