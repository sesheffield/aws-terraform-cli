package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var cmdS3Bucket = &cobra.Command{
	Use:   "s3 <sub-command>",
	Short: "Terraform S3",
	Long:  `S3 is for creating terraform configurations.`,
	//Run:   S3Run,
}

var cmdCreateS3Bucket = &cobra.Command{
	Use:   "create [s3-bucket-name]",
	Short: "Creating a S3 Bucket",
	Long:  `Create is for terraforming a S3 Bucket.`,
	Run:   createS3Bucket,
}

var cmdS3BucketEvent = &cobra.Command{
	Use:   "event [service]",
	Short: "Event will connect a service to a S3 Bucket",
	Long:  `Event will terraform a notification to a S3 bucket based on a lambda function`,
	Run:   createS3Bucket,
}

var (
	tagName              string
	tagProject           string
	acl                  string
	serverSideEncryption string
	environment          string
	author               string
)

func init() {
	rootCmd.AddCommand(cmdS3Bucket)
	cmdS3Bucket.AddCommand(cmdCreateS3Bucket)
	cmdCreateS3Bucket.Flags().StringVarP(&environment, "env", "e", "", "Environment")
}

// Need to create s3 bucket with lambda permisions,
// potential notification triggers,
// iam access
// Needs to be able to check if s3 bucket already exists and if it does then to add to it (for notifications or permisions)
func createS3Bucket(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		panic(errors.New("missing s3 bucket name"))
	}
	if environment == "" {
		fmt.Println("Using dev environment by default")
		environment = "dev"
	}

	// If the file does not exist, create it, or append to the file
	s3File, err := os.OpenFile("s3.tf", os.O_APPEND|os.O_CREATE|os.O_RDWR, 7777)
	if err != nil {
		panic(err)
	}
	defer s3File.Close()

	bucketName := args[0]
	checkS3Bucket := fmt.Sprintf("resource \"aws_s3_bucket\" \"%s-%s\"", bucketName, environment)
	reader := bufio.NewReader(s3File)
	var duplicate bool
	for {
		line, err := reader.ReadString('\n')
		if strings.Contains(line, checkS3Bucket) {
			duplicate = true
			break
		}
		if err == io.EOF {
			break
		}
	}

	if !duplicate {
		passedVars := Vars{
			Required: Required{
				S3BucketName: bucketName,
				Environment:  environment,
			},
			Optional: Optional{
				ACL:                  aclTF,
				ServerSideEncryption: serverSideEncryptionTF,
			},
		}

		s3BucketTemplate := template.Must(template.New("s3Bucket").Parse(s3BucketTF))
		var s3BucketTPL bytes.Buffer
		if err := s3BucketTemplate.Execute(&s3BucketTPL, &passedVars); err != nil {
			panic(err)
		}
		s3TF := strings.Replace(s3BucketTPL.String(), "\n\t\n", "\n", -1)
		_, err = s3File.WriteString(s3TF)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("%s already exists for your %s environment\n", bucketName, environment)
	}
}

const s3BucketTF = `
resource "aws_s3_bucket" "{{.Required.S3BucketName}}-{{.Required.Environment}}" {
	bucket = "{{.Required.S3BucketName}}-{{.Required.Environment}}"
	{{.Optional.ACL}}
	{{.Optional.ServerSideEncryption}}
	tags {
		Name        = "${vars.project_name}-{{.Required.Environment}}-tfstate-${vars.concat_aws_region}"
    Project     = "${vars.project_name}"
    Environment = "{{.Required.Environment}}"
    Author      = "Terraform"
  }
}
`
