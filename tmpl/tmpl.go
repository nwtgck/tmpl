package tmpl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type TmplYaml struct {
	// NOTE: yaml.MapSlice preserves the order of items
	// (from: https://stackoverflow.com/a/42109240/2885946)
	Variables yaml.MapSlice `yaml:variables`
}

const TmplYamlName = "tmpl.yaml"

func FillVariables(dirPath string) error {
	tmplYaml, err := ReadTemplYaml(dirPath)
	if err != nil {
		return err
	}
	// Input variable values from user input
	inputVariables := InputVariables(tmplYaml.Variables)

	// Combine reserved variables
	variables := GetReservedVariables()
	for name, value := range inputVariables {
		variables[name] = value
	}

	// Replace files in the directory
	err = ReplaceInDir(dirPath, variables)
	return err
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

func getDiffs(dmp *diffmatchpatch.DiffMatchPatch, original string, filled string) []diffmatchpatch.Diff {
	// Calculate diffs between original and filled one
	diffs := dmp.DiffMain(original, filled, false)
	// Compact diffs
	compactDiffs := getCompactDiffs(diffs)

	return compactDiffs
}

func ReplaceInDir(dirPath string, variables map[string]interface{}) error {
	dmp := diffmatchpatch.New()
	// Each file in the root directory
	// (from: https://flaviocopes.com/go-list-files/)
	err := filepath.Walk(dirPath, func(fpath string, info os.FileInfo, err error) error {
		// If fpath is in .git directory
		// TODO: Use better way
		if strings.Contains(fpath, ".git/") {
			// Skip
			return nil
		// If fpath is not directory fpath
		} else if !info.IsDir(){
			// Read whole file content
			original, _ := ioutil.ReadFile(fpath)
			// (from: https://stackoverflow.com/a/49043639/2885946)
			name := path.Base(fpath)
			// Create a new template and parse the letter into it.
			// (from: https://golang.org/pkg/text/template/#example_Template)
			t, err := template.New(name).ParseFiles(fpath)
			if err != nil {
				return err
			}
			buf := &bytes.Buffer{}
			err = t.Execute(buf, variables)
			if err != nil {
				return err
			}
			// Overwrite filled one
			ioutil.WriteFile(fpath, buf.Bytes(), info.Mode())

			// Get diffs
			diffs := getDiffs(dmp, string(original), buf.String())
			// If there are diffs
			if len(diffs) != 0 {
				// Print diffs
				fmt.Printf("====== %s ======\n", fpath)
				fmt.Println(dmp.DiffPrettyText(diffs))
			}
		}
		return nil
	})
	return err
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

func GetReservedVariables() map[string]interface{} {
	return map[string]interface{} {
		"Now": map[string]interface{}{
			"year": time.Now().Year(),
		},
		"Git": map[string]interface{}{
			"user": map[string]string{
				"name": gitUserName(),
				"email": gitUserEmail(),
			},
		},
	}
}