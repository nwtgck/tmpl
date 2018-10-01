package cmd

import (
	"fmt"
	"github.com/nwtgck/tmpl/tmpl"
	"github.com/spf13/cobra"
	"os"
)

// Dry-run fill
var fillDryRun bool
// Fill by YAML
var fillFillYamlStr string

func init() {
	RootCmd.AddCommand(fillCmd)
	fillCmd.Flags().BoolVarP(&fillDryRun, "dry-run", "n", false, "dry-run")
	fillCmd.Flags().StringVar(&fillFillYamlStr, "fill-yaml", "", "fill variables by YAML")
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
		err := tmpl.FillVariables(dirPath, fillFillYamlStr, true, fillDryRun) // TODO: enableYamlParse is hard coded
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Fill failed with %s.\n", err)
			os.Exit(-1)
		}
	},
}
