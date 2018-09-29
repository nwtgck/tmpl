package main

import (
	"bufio"
	"fmt"
	"github.com/cbroglie/mustache"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TmplYaml struct {
	Variables map[string]string `yaml:variables`
}

const TmplYamlName = "tmpl.yaml"

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

func readTemplYaml(dirPath string) (TmplYaml, error) {
	// TODO: Check existence of tmpl.yaml
	buf, err := ioutil.ReadFile(path.Join(dirPath, TmplYamlName))
	// Create
	var tmplYaml TmplYaml
	err = yaml.Unmarshal(buf, &tmplYaml)
	if err != nil {
		return tmplYaml, err
	}
	return tmplYaml, err
}

func gitUserName() string {
	out, err := exec.Command("git", "config", "user.name").Output()
	if err != nil || len(out) == 0 {
		return ""
	} else {
		return strings.TrimRight(string(out), "\n")
	}
}

func gitUserEmail() string {
	out, err := exec.Command("git", "config", "user.email").Output()
	if err != nil || len(out) == 0 {
		return ""
	} else {
		return strings.TrimRight(string(out), "\n")
	}
}

func getReservedVariables() map[string]string {
	return map[string]string {
		"$year": strconv.Itoa(time.Now().Year()),
		"$git_user_name": gitUserName(),
		"$git_user_email": gitUserEmail(),
	}
}

func main(){
	// Get root directory path
	// TODO: Error handling
	dirPath := os.Args[1]
	tmplYaml, err := readTemplYaml(dirPath)
	if err != nil {
		panic(err)
	}
	// Input variable values from user input
	variables := inputVariables(tmplYaml.Variables)

	// Combine reserved variables
	for name, value := range getReservedVariables() {
		variables[name] = value
	}

	// Replace files in the directory
	replaceInDir(dirPath, variables)
}
