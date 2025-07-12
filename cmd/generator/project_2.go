package generator

import (
	projectv2 "github.com/AugustineAurelius/eos/generator/project-v2"
	"github.com/spf13/cobra"
)

func projectV2CMD() *cobra.Command {

	var pathToSpec string
	cmd := &cobra.Command{
		Use:   "project-v2",
		Short: "project v2 generator",

		RunE: func(cmd *cobra.Command, args []string) error {

			return projectv2.Generate(pathToSpec)
		},
	}

	cmd.PersistentFlags().StringVarP(&pathToSpec, "spec", "s", "", "path to spec")

	cmd.MarkPersistentFlagRequired("spec")

	return cmd
}
