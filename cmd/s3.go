package cmd

import "github.com/spf13/cobra"

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

var (
	tagName    string
	tagProject string
)

func init() {
	cmdS3Bucket.AddCommand(cmdCreateS3Bucket)
	//cmdCreateS3Bucket.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")

}

// Need to create s3 bucket with lambda permisions,
// potential notification triggers,
// iam access
// Needs to be able to check if s3 bucket already exists and if it does then to add to it (for notifications or permisions)
func createS3Bucket(cmd *cobra.Command, args []string) {

}
