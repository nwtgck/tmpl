package cmd

import (
	"github.com/nwtgck/tmpl/version"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:     os.Args[0],
	Short:   "Trans CLI",
	Long:    "Trans CLI",
	Version: version.Version,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
