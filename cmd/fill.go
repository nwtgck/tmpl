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
		// TODO: Error handling
		dirPath := args[0]
		tmplYaml, err := tmpl.ReadTemplYaml(dirPath)
		if err != nil {
			panic(err)
		}
		// Input variable values from user input
		variables := tmpl.InputVariables(tmplYaml.Variables)

		// Combine reserved variables
		for name, value := range tmpl.GetReservedVariables() {
			variables[name] = value
		}

		// Replace files in the directory
		tmpl.ReplaceInDir(dirPath, variables)
	},
}
