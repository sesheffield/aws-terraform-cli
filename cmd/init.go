package cmd

import (
	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:   "init <sub-command>",
	Short: "Initilize Terraform Environment",
	Long:  `init will create the backend.tf and vars-config.tf files.`,
	Run:   initRun,
}

func init() {
	rootCmd.AddCommand(cmdInit)
}

func initRun(cmd *cobra.Command, args []string) {
	init := NewInitFacade()
	init.start()
}
