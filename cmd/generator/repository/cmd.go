package repository

import (
	"fmt"

	repositorygen "github.com/AugustineAurelius/eos/repository_gen"
	"github.com/spf13/cobra"
)

var (
	StructName string
)

var Cmd = &cobra.Command{
	Use:   "repository",
	Short: "repository pattern",
	Long:  `repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to generate builder pattern for struct")

		repositorygen.Generate(StructName)
	},
}
