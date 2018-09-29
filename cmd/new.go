package cmd

import (
	"fmt"
	"github.com/nwtgck/tmpl/tmpl"
	"github.com/spf13/cobra"
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
		if (len(args) == 0) {
			fmt.Fprintln(os.Stderr, "Error: Git repository path is not specified")
			os.Exit(-1)
		}
		// Get git repository path
		gitRepoPath := args[0]
		// Git clone
		fmt.Printf("Cloning '%s'...\n", gitRepoPath)
		err := exec.Command("git", "clone", "--recursive", gitRepoPath).Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Clone failed with %s\n", err)
			os.Exit(-1)
		}
		// Git dir path
		dirPath := getFileNameWithoutExt(gitRepoPath)
		// Fill .tmpl with variables
		tmpl.FillVariables(dirPath)
	},
}

// (from: https://qiita.com/KemoKemo/items/d135ddc93e6f87008521)
func getFileNameWithoutExt(path string) string {
	// Fixed with a nice method given by mattn-san
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}