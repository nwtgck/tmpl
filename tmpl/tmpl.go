package tmpl

import (
	"bufio"
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/sergi/go-diff/diffmatchpatch"
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
	// NOTE: yaml.MapSlice preserves the order of items
	// (from: https://stackoverflow.com/a/42109240/2885946)
	Variables yaml.MapSlice `yaml:variables`
}

const TmplYamlName = "tmpl.yaml"

func FillVariables(dirPath string) {
	tmplYaml, err := ReadTemplYaml(dirPath)
	if err != nil {
		panic(err)
	}
	// Input variable values from user input
	variables := InputVariables(tmplYaml.Variables)

	// Combine reserved variables
	for name, value := range GetReservedVariables() {
		variables[name] = value
	}

	// Replace files in the directory
	ReplaceInDir(dirPath, variables)
}

func getCompactDiffs(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := []diffmatchpatch.Diff{}
	hasDiffInLine := false
	lineDiffs := []diffmatchpatch.Diff{}
	for _, diff := range diffs {
		hasDiff := diff.Type != diffmatchpatch.DiffEqual
		hasNewLine := strings.Contains(diff.Text, "\n")
		if hasDiff {
			// Set flag
			hasDiffInLine = true
		}
		// If diff has newline
		if hasNewLine {
			// If diffs have some diff in line
			if hasDiffInLine {
				result = append(result, lineDiffs...)
			}
			hasDiffInLine = false
			lineDiffs = []diffmatchpatch.Diff{}
		}

		// Append to line-diffs
		lineDiffs = append(lineDiffs, diff)
	}
	// If diffs have some diff in line
	if hasDiffInLine {
		result = append(result, lineDiffs...)
	}
	return result
}

func ReplaceInDir(dirPath string, variables map[string]string) error {
	dmp := diffmatchpatch.New()
	// Each file in the root directory
	// (from: https://flaviocopes.com/go-list-files/)
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// If path is not directory path
		if !info.IsDir(){
			// Read whole file content
			original, _ := ioutil.ReadFile(path)
			// Filled with variables
			data, _  := mustache.RenderFile(path, variables)
			// Overwrite filled one
			ioutil.WriteFile(path, []byte(data), info.Mode())
			// Calculate diffs between original and filled one
			diffs := dmp.DiffMain(string(original), data, false)
			// Compact diffs
			compactDiffs := getCompactDiffs(diffs)
			// If there are diffs
			if len(compactDiffs) != 0 {
				// Print diffs
				fmt.Printf("====== %s ======\n", path)
				fmt.Println(dmp.DiffPrettyText(compactDiffs))
			}
		}
		return nil
	})
	return nil
}

func InputVariables(prompt yaml.MapSlice) map[string]string {
	scanner := bufio.NewScanner(os.Stdin)
	variables := map[string]string{}
	for _, item := range prompt {
		// Get variable name
		varName := item.Key.(string)
		// Get description
		desc    := item.Value.(string)
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

func ReadTemplYaml(dirPath string) (TmplYaml, error) {
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

func GetReservedVariables() map[string]string {
	return map[string]string {
		"$year": strconv.Itoa(time.Now().Year()),
		"$git_user_name": gitUserName(),
		"$git_user_email": gitUserEmail(),
	}
}