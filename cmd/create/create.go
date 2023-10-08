package create

import (
	"fmt"
	"main/scaffold/create"

	"github.com/spf13/cobra"
)

var (
	dryRun     bool
	production bool
	yamlF      string
	jsonF      string
	Cmd        = &cobra.Command{
		Use:  "create <environment-name> -y <yaml-filename> -j <json-filename> [-d]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if dryRun {
				fmt.Println("Running dryrun mode")
			}
			if yamlF != "" {
				create.CreateScaffold(args[0], production)
			}
		},
	}
)

func init() {
	Cmd.Flags().BoolVarP(&dryRun, "dryrun", "d", false, "outputs to stdout")
	Cmd.Flags().BoolVarP(&production, "production", "p", false, "sets NODE_ENV property to Production if true")
	Cmd.Flags().StringVarP(&yamlF, "yaml", "y", "", "yaml filename to scaffold")
	Cmd.Flags().StringVarP(&jsonF, "json", "j", "", "json filename to scaffold")
}
