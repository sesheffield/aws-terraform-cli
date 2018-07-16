package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var cmdEcho = &cobra.Command{
	Use:   "echo [string to echo]",
	Short: "Echo anything to the screen",
	Long: `echo is for echoing anything back.
    Echo echo’s.
    `,
	Run: echoRun,
}

var times int
var times2 int

func init() {
	rootCmd.AddCommand(cmdEcho)
	cmdEcho.AddCommand(versionCmd)
	cmdEcho.Flags().IntVarP(&times, "times", "n", 1, "times to echo")
	cmdEcho.Flags().IntVarP(&times2, "times2", "t", 1, "times to echo again")
}

func echoRun(cmd *cobra.Command, args []string) {
	for i := 0; i < times; i++ {
		fmt.Println(strings.Join(args, " "))
	}
	for i := 0; i < times2; i++ {
		fmt.Println(strings.Join(args, " "))
	}
}

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		versionRun(cmd, args)
	},
}
