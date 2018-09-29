package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:     os.Args[0],
	Short:   "Trans CLI",
	Long:    "Trans CLI",
	Version: "dummy-version", // TODO: Fill version
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
