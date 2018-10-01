package cmd

import (
	"fmt"
	"github.com/nwtgck/tmpl/tmpl"
	"github.com/nwtgck/tmpl/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
)

// Dry-run fill
var newDryRun bool
// Fill by YAML
var newFillYamlStr string

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&newDryRun, "dry-run", "n", false, "dry-run")
	newCmd.Flags().StringVar(&newFillYamlStr, "fill-yaml", "", "fill variables by YAML")
}

var newCmd = &cobra.Command{
	Use: "new",
	Short: "Download template from Git repo and Fill",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: Git repository path is not specified")
			os.Exit(-1)
		}
		// Get git repository path
		gitRepoPath := args[0]
		// Git dir path
		dirPath := util.GetFileNameWithoutExt(gitRepoPath)
		if len(args) == 2 {
			dirPath = args[1]
		}
		if !newDryRun && util.Exists(dirPath) {
			if util.Ask4confirm(fmt.Sprintf("Are you sure to overwrite '%s'? ", dirPath)) {
				os.RemoveAll(dirPath)
			} else {
				fmt.Println("Canceled.")
				os.Exit(0)
			}
		}
		tmpRepoPath, err := ioutil.TempDir("", "repo")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Temporary directory creation failed with %s\n", err)
			os.Exit(-1)
		}
		// Git clone
		fmt.Printf("Cloning '%s'...\n", gitRepoPath)
		err = exec.Command("git", "clone", "--recursive", gitRepoPath, tmpRepoPath).Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Clone failed with %s\n", err)
			os.Exit(-1)
		}
		if newDryRun {
			// Use temporary repo path
			dirPath = tmpRepoPath
		} else {
			// Move tmpRepoPath to dirPath
			err = os.Rename(tmpRepoPath, dirPath)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Rename failed with %s\n", err)
			os.Exit(-1)
		}

		// Fill .tmpl with variables
		// TODO: Hard code fill yaml
		err = tmpl.FillVariables(dirPath, newFillYamlStr, true, newDryRun) // TODO: enableYamlParse is hard coded
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Fill failed with %s.\n", err)
			os.Exit(-1)
		}
		// Clean up temp directory
		defer os.RemoveAll(tmpRepoPath)
	},
}
