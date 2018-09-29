package cmd

import (
	"github.com/nwtgck/tmpl/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(fillCmd)
}

var fillCmd = &cobra.Command{
	Use: "fill",
	Short: "Fill .tmpl/ with variables",
	Run: func(cmd *cobra.Command, args []string) {
		// Get root directory path
		// (from: https://stackoverflow.com/a/31483763/2885946)
		dirPath := "."
		if len(args) >= 1 {
			dirPath = args[0]
		}
		tmpl.FillVariables(dirPath)
	},
}
