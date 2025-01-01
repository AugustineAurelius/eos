package repository

import (
	"fmt"

	repositorygen "github.com/AugustineAurelius/eos/generator/repository"
	"github.com/spf13/cobra"
)

var (
	StructName, TxRunenrPath string
	WithTX                   bool
)

var Cmd = &cobra.Command{
	Use:   "repository",
	Short: "repository pattern",
	Long:  `repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to generate repository for struct")

		repositorygen.Generate(StructName, TxRunenrPath, WithTX)
	},
}
