package generator

import (
	"fmt"

	repositorygen "github.com/AugustineAurelius/eos/generator/repository"

	"github.com/spf13/cobra"
)

func repositoryCMD() *cobra.Command {
	var (
		structName    string
		WithDefaultID bool
	)
	cmd := &cobra.Command{
		Use:   "repository",
		Short: "repository pattern",
		Long:  `repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start to generate repository for struct")

			repositorygen.Generate(structName, WithDefaultID)
		},
	}

	cmd.PersistentFlags().StringVarP(&structName, "type", "t", "", "name of the struct for which would be generated repo")
	cmd.PersistentFlags().BoolVarP(&WithDefaultID, "default_id", "i", false, "add id to create")
	return cmd
}
