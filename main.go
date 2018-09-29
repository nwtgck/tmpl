package main

import (
	"bufio"
	"fmt"
	"github.com/cbroglie/mustache"
	"io/ioutil"
	"os"
	"path/filepath"
)


func replaceInDir(dirPath string, variables map[string]string) error {
	// Each file in the root directory
	// (from: https://flaviocopes.com/go-list-files/)
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// If path is not directory path
		if !info.IsDir(){
			// Just print file path
			fmt.Printf("====== %s ======\n", path)
			data, _  := mustache.RenderFile(path, variables)
			ioutil.WriteFile(path, []byte(data), info.Mode())
		}
		return nil
	})
	return nil
}

func inputVariables(prompt map[string]string) map[string]string {
	scanner := bufio.NewScanner(os.Stdin)
	variables := map[string]string{}
	for varName, desc := range prompt {
		// Print variable name and description
		fmt.Printf("%s (%s) = ", varName, desc)
		// Get line
		scanner.Scan()
		line := scanner.Text()
		// Add pair of variable name and its value
		variables[varName] = line
	}
	return variables
}

func main(){
	// TODO: Hard code (this will be defined by .yaml)
	prompt := map[string]string{
		"myvar1": "this is myvar1",
		"myvar2": "this is not myvar1 but 2",
	}
	// Input variable values from user input
	variables := inputVariables(prompt)
	// Get root directory path
	// TODO: Error handling
	dirPath := os.Args[1]
	// Replace files in the directory
	replaceInDir(dirPath, variables)
}
