package tmpl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nwtgck/tmpl/util"
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

func FillVariables(dirPath string, enableYamlParse bool, dryRun bool) error {
	tmplYaml, err := ReadTemplYaml(dirPath)
	if err != nil {
		return err
	}
	// Input variable values from user input
	inputVariables := InputVariables(tmplYaml.Variables, enableYamlParse)

	// Combine reserved variables
	variables := GetReservedVariables(enableYamlParse)
	for name, value := range inputVariables {
		variables[name] = value
	}

	// Replace files in the directory
	err = ReplaceInDir(dirPath, variables, dryRun)
	return err
}

func getDiffs(dmp *diffmatchpatch.DiffMatchPatch, original string, filled string) []diffmatchpatch.Diff {
	// Calculate diffs between original and filled one
	// (from: https://qiita.com/shibukawa/items/dd75ad01e623c4c1166b)
	a, b, c := dmp.DiffLinesToChars(original, filled)
	diffs := dmp.DiffMain(a, b, false)
	lineBasedDiffs := dmp.DiffCharsToLines(diffs, c)

	// Use only not-equal diff
	result := []diffmatchpatch.Diff{}
	for _, d := range lineBasedDiffs {
		if d.Type != diffmatchpatch.DiffEqual {
			// Prepend "+" / "-"
			switch d.Type {
			case diffmatchpatch.DiffInsert:
				d.Text = "+ " + d.Text
			case diffmatchpatch.DiffDelete:
				d.Text = "- " + d.Text
			}
			result = append(result, d)
		}
	}
	return result
}

func ReplaceInDir(dirPath string, variables map[string]interface{}, dryRun bool) error {
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
			// Not dry-run
			if !dryRun {
				// Overwrite filled one
				ioutil.WriteFile(fpath, buf.Bytes(), info.Mode())
			}

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

func InputVariables(prompt yaml.MapSlice, enableYamlParse bool) map[string]interface{} {
	scanner := bufio.NewScanner(os.Stdin)
	variables := map[string]interface{}{}
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

		var value interface{}
		if enableYamlParse {
			// Parse line and assign into value
			err := yaml.Unmarshal([]byte(line), &value)
			if err != nil {
				value = line
			}
			//fmt.Println(value)
		} else {
			value = line
		}
		// Add pair of variable name and its value
		variables[varName] = value
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

func GetReservedVariables(enableYamlParse bool) map[string]interface{} {
	env := map[string]interface{}{}
	if enableYamlParse {
		// Parse environment variable
		for name, valueStr := range util.GetEnv() {
			var value interface{}
			// Parse value string
			err := yaml.Unmarshal([]byte(valueStr), &value)
			if err == nil {
				env[name] = value
			} else {
				env[name] = valueStr
			}
		}
	} else {
		for name, valueStr := range util.GetEnv() {
			env[name] = valueStr
		}
	}

	return map[string]interface{} {
		"Env": env,
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