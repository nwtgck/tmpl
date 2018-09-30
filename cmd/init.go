package cmd

import (
	"fmt"
	"github.com/nwtgck/tmpl/tmpl"
	"github.com/nwtgck/tmpl/util"
	"github.com/spf13/cobra"
	"github.com/MakeNowJust/heredoc"
	"io/ioutil"
	"os"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initYaml = heredoc.Doc(fmt.Sprintf(`
# This is an example of %s

# (Structure is <variable name>: <description>)
# variables:
#  project_name: Project Name
`, tmpl.TmplYamlName))

var initCmd = &cobra.Command{
	Use: "init",
	Short: fmt.Sprintf("Create %s", tmpl.TmplYamlName),
	Run: func(cmd *cobra.Command, args []string) {
		if util.Exists(tmpl.TmplYamlName) {
			if !util.Ask4confirm(fmt.Sprintf("Are you sure to overwrite '%s'? ", tmpl.TmplYamlName)) {
				fmt.Println("Canceled.")
				os.Exit(0)
			}
		}
		err := ioutil.WriteFile(tmpl.TmplYamlName, []byte(initYaml), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Write failed with %s\n", err)
			os.Exit(-1)
		}
		fmt.Printf("'%s' created!\n", tmpl.TmplYamlName)
	},
}
