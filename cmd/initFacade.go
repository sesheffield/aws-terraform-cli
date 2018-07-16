package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
)

const stateS3Bucket = `
# 
# Bucket for Terraform state files
#
resource "aws_s3_bucket" "tfstate_%s" {
  bucket = "${var.project_name}-tfstate-%s"
  acl    = "private"
  versioning {
    enabled = true
  }
  tags {
    Name    = "${var.project_name}-tfstate-%s"
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
	bucket = "%s-tfstate-%s"
	region = "%s"
	key = "%s-apps-ops.tfstate"
	encrypt = "true"
	acl = "private"
  }
}
`

const awsProviderTF = `
# AWS PROVIDER
provider "aws" {
  region = "${var.aws_region}"
  allowed_account_ids = [
    "${var.aws_account_id}",
  ]
  version = "~> 1.3"
}
`
const awsAccountIDTF = `
# AWS ACCOUNT ID
variable "aws_account_id" {
  default = "%d"
}
`
const awsRegionTF = `
# AWS REGION
variable "aws_region" {
  default = "%s"
}
`
const projectNameTF = `
# AWS PROJECT NAME 
variable "project_name" {
  default = "%s"
}
`

const provisionerTF = `
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

type varsTerraform struct {
	awsAccountID int64
	awsRegion    string
	projectName  string
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

func (t *initTerraform) terraformInit(vars *varsTerraform) error {
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

	awsRegionRaw := strings.Replace(vars.awsRegion, "-", "", -1)
	coreStateTF := fmt.Sprintf(coreState, vars.projectName, awsRegionRaw, vars.awsRegion, vars.projectName)
	varsTF := awsProviderTF +
		fmt.Sprintf(awsAccountIDTF, vars.awsAccountID) +
		fmt.Sprintf(awsRegionTF, vars.awsRegion) +
		fmt.Sprintf(projectNameTF, vars.projectName) +
		provisionerTF
	s3TF := fmt.Sprintf(stateS3Bucket, awsRegionRaw, awsRegionRaw, awsRegionRaw)

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

func (u *userInput) handleUserInput(stdin *bufio.Reader) (*varsTerraform, error) {

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
		return nil, err
	}
	projectName = strings.TrimRightFunc(projectName, func(c rune) bool {
		//In windows newline is \r\n
		return c == '\r' || c == '\n'
	})
	return &varsTerraform{
		awsAccountID: awsAccountID,
		awsRegion:    awsRegion,
		projectName:  projectName,
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
