package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/apex/log"
)

const stateS3Bucket = `
#
# Bucket for Terraform state files
#
resource "aws_s3_bucket" "{{.Required.StateBucket}}" {
  bucket = "${var.project_name}-{{.Required.StateBucket}}-{{.Required.ConcatRegion}}"
	{{.Optional.ACL}}
	{{.Optional.Versioning}}
	tags {
		Name    = "${var.project_name}-{{.Required.StateBucket}}-{{.Required.ConcatRegion}}"
    Project = "${var.project_name}"
    Meta    = "author:${var.provisioner}"
  }
}
`

const coreState = `#
# Core state
#
terraform {
	backend "s3" {
	bucket = "{{.Required.ProjectName}}-{{.Required.StateBucket}}-{{.Required.Region}}"
	region = "{{.Required.Region}}"
	key = "{{.Required.ProjectName}}-apps-ops.{{.Required.StateBucket}}"
	{{.Optional.Encrypt}}
	{{.Optional.ACL}}
  }
}
`

const defaultVars = `
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
  default = "{{.Required.AccountID}}"
}

# AWS REGION
variable "aws_region" {
  default = "{{.Required.Region}}"
}

# CONCAT AWS REGION
variable "concat_aws_region" {
  default = "{{.Required.ConcatRegion}}"
}

# AWS PROJECT NAME 
variable "project_name" {
  default = "{{.Required.ProjectName}}"
}

# AWS PROVISIONER
variable "provisioner" {
  default = "terraform"
}
`

var validAWSRegions = map[string]struct{}{
	"us-east-1":      struct{}{},
	"us-east-2":      struct{}{},
	"us-west-1":      struct{}{},
	"us-west-2":      struct{}{},
	"ca-central-1":   struct{}{},
	"eu-central-1":   struct{}{},
	"eu-west-1":      struct{}{},
	"eu-west-2":      struct{}{},
	"eu-west-3":      struct{}{},
	"ap-northeast-1": struct{}{},
	"ap-northeast-2": struct{}{},
	"ap-northeast-3": struct{}{},
	"ap-southeast-1": struct{}{},
	"ap-southeast-2": struct{}{},
	"ap-south-1":     struct{}{},
	"sa-east-1":      struct{}{},
}

type userInput struct{}
type initTerraform struct{}

// InitFacade ...
type InitFacade struct {
	userInput       *userInput
	createInitFiles *initTerraform
}

// NewInitFacade ...
func NewInitFacade() *InitFacade {
	return &InitFacade{new(userInput), new(initTerraform)}
}

func (c *InitFacade) start() {
	stdin := bufio.NewReader(os.Stdin)
	results, err := c.userInput.handleUserInput(stdin)
	if err != nil {
		panic(err)
	}
	if err := c.createInitFiles.terraformInit(results); err != nil {
		panic(err)
	}
}

func (t *initTerraform) terraformInit(vars Vars) error {
	varsConfigFile, err := os.Create("vars-config.tf")
	if err != nil {
		log.WithError(err).Error("unable to create vars config file")
		return err
	}
	defer varsConfigFile.Close()

	backEndFile, err := os.Create("backend.tf")
	if err != nil {
		log.WithError(err).Error("unable to create backend config file")
		return err
	}
	defer backEndFile.Close()

	s3StateFile, err := os.Create("s3.tf")
	if err != nil {
		log.WithError(err).Error("unable to create s3 state file")
		return err
	}
	defer s3StateFile.Close()

	vars.ConcatRegion = strings.Replace(vars.Required.Region, "-", "", -1)

	coreStateTemplate := template.Must(template.New("coreState").Parse(coreState))
	var coreTPL bytes.Buffer
	if err := coreStateTemplate.Execute(&coreTPL, &vars); err != nil {
		panic(err)
	}
	coreStateTF := strings.Replace(coreTPL.String(), "\n\t\n", "\n", -1)

	defaultVarsTemplate := template.Must(template.New("defaultVars").Parse(defaultVars))
	var defaultVarsTPL bytes.Buffer
	if err := defaultVarsTemplate.Execute(&defaultVarsTPL, &vars); err != nil {
		panic(err)
	}
	varsTF := strings.Replace(defaultVarsTPL.String(), "\n\t\n", "\n", -1)

	stateS3BucketTemplate := template.Must(template.New("stateS3Bucket").Parse(stateS3Bucket))
	var stateTPL bytes.Buffer

	if err := stateS3BucketTemplate.Execute(&stateTPL, &vars); err != nil {
		panic(err)
	}

	s3TF := strings.Replace(stateTPL.String(), "\n\t\n", "\n", -1)

	_, err = varsConfigFile.WriteString(varsTF)
	if err != nil {
		panic(err)
	}

	_, err = backEndFile.WriteString(coreStateTF)
	if err != nil {
		panic(err)
	}

	_, err = s3StateFile.WriteString(s3TF)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nCreated backend.tf, vars-config.tf, and s3.tf files")
	return nil
}

func (u *userInput) handleUserInput(stdin *bufio.Reader) (Vars, error) {

	fmt.Print("AWS Account ID: ")
	awsAccountID := readIntInput(stdin, "AWS Account ID")

	fmt.Print("AWS Region: ")
	awsRegion, err := handleAWSRegion(stdin)
	if err != nil {
		panic(err)
	}

	fmt.Print("Project Name: ")
	projectName, err := stdin.ReadString('\n')
	if err != nil {
		log.WithError(err).Error("unable to read project name input")
		return Vars{}, err
	}
	projectName = strings.TrimRightFunc(projectName, func(c rune) bool {
		//In windows newline is \r\n
		return c == '\r' || c == '\n'
	})

	return Vars{
		Required: Required{
			StateBucket: "tfstate",
			AccountID:   awsAccountID,
			Region:      awsRegion,
			ProjectName: projectName,
		},
		Optional: Optional{
			ACL:        aclTF,
			Versioning: versioningTF,
			Encrypt:    encryptTF,
		},
	}, nil
}

func readIntInput(stdin *bufio.Reader, name string) int64 {
	var val int64
	for {
		_, err := fmt.Fscan(stdin, &val)
		if err == nil {
			break
		}
		// will only print below once and wait for the newline
		stdin.ReadString('\n')
		fmt.Print(fmt.Sprintf("Sorry, invalid %s provided. Please enter an integer: ", name))
	}
	// clearing from last newline
	stdin.ReadString('\n')
	return val
}

func handleAWSRegion(stdin *bufio.Reader) (string, error) {
	var (
		awsRegion string
		err       error
	)
	for {
		awsRegion, err = stdin.ReadString('\n')
		if err != nil {
			log.WithError(err).Error("unable to read aws region input")
			return "", err
		}
		awsRegion = strings.TrimRightFunc(awsRegion, func(c rune) bool {
			//In windows newline is \r\n
			return c == '\r' || c == '\n'
		})
		// check that aws region is valid
		if _, ok := validAWSRegions[awsRegion]; ok {
			break
		}
		// will only print below once and wait for the newline
		fmt.Print("Sorry, invalid region provided. Please enter a valid region: ")
	}
	return awsRegion, nil
}
