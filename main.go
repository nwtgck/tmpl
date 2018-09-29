package main

import (
	"os"
	"github.com/nwtgck/tmpl/tmpl"
)

func main(){
	// Get root directory path
	// TODO: Error handling
	dirPath := os.Args[1]
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
}
