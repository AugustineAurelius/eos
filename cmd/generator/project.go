package generator

import (
	"github.com/AugustineAurelius/eos/generator/project"
	"github.com/spf13/cobra"
)

func projectCMD() *cobra.Command {

	var outputDir, projectName, url string
	cmd := &cobra.Command{
		Use:   "project",
		Short: "project generator",

		RunE: func(cmd *cobra.Command, args []string) error {

			return project.Generate(project.ProjectData{
				ProjectURL:  url,
				ProjectName: projectName,
				Output:      outputDir,
			})
		},
	}

	cmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "output dir path")
	cmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "name of project")
	cmd.PersistentFlags().StringVarP(&url, "url", "u", "", "path to repos")

	cmd.MarkPersistentFlagRequired("project")
	cmd.MarkPersistentFlagRequired("github")

	return cmd
}
