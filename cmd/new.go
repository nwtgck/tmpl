package cmd

import (
	"fmt"
	"github.com/nwtgck/tmpl/tmpl"
	"github.com/nwtgck/tmpl/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	RootCmd.AddCommand(newCmd)
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
		// Git tmpRepoPath path
		dirPath := getFileNameWithoutExt(gitRepoPath)
		if len(args) == 2 {
			dirPath = args[1]
		}
		if util.Exists(dirPath) {
			fmt.Fprintf(os.Stderr, "Error: '%s' already exists.\n", dirPath)
			os.Exit(-1)
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
		// Move tmpRepoPath to dirPath
		err = os.Rename(tmpRepoPath, dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Rename failed with %s\n", err)
			os.Exit(-1)
		}
		// Clean up temp directory
		defer os.RemoveAll(tmpRepoPath)

		// Fill .tmpl with variables
		err = tmpl.FillVariables(dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Fill failed with %s.\n", err)
			os.Exit(-1)
		}
	},
}

// (from: https://qiita.com/KemoKemo/items/d135ddc93e6f87008521)
func getFileNameWithoutExt(path string) string {
	// Fixed with a nice method given by mattn-san
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}