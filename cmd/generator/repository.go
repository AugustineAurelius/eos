package generator

import (
	"fmt"

	repositorygen "github.com/AugustineAurelius/eos/generator/repository"

	"github.com/spf13/cobra"
)

func repositoryCMD() *cobra.Command {
	var (
		structName, txRunenrPath, commonPath string
		withTX                               bool
	)
	cmd := &cobra.Command{
		Use:   "repository",
		Short: "repository pattern",
		Long:  `repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start to generate repository for struct")

			repositorygen.Generate(structName, txRunenrPath, commonPath, withTX)
		},
	}

	cmd.PersistentFlags().StringVarP(&structName, "type", "t", "", "name of the struct for which would be generated repo")
	cmd.PersistentFlags().BoolVarP(&withTX, "with_tx", "x", false, "does repo would be with transaction logic")
	cmd.PersistentFlags().StringVarP(&txRunenrPath, "tx_path", "p", "", "path to txrunner pkg")
	cmd.PersistentFlags().StringVarP(&commonPath, "common_path", "c", "", "path to common pkg")
	cmd.MarkPersistentFlagRequired("common_path")

	return cmd
}
