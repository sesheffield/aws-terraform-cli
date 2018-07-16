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
	//cmdInit.AddCommand(versionCmd)
	//cmdInit.Flags().IntVarP(&times, "times", "n", 1, "times to echo")
	//cmdInit.Flags().IntVarP(&times2, "times2", "t", 1, "times to echo again")
}

func initRun(cmd *cobra.Command, args []string) {
	init := NewInitFacade()
	init.start()
}
