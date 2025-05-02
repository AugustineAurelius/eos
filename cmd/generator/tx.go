package generator

import (
	"fmt"

	txrunner "github.com/AugustineAurelius/eos/generator/tx_runner"

	"github.com/spf13/cobra"
)

func txCMD() *cobra.Command {
	var commonPath string

	cmd := &cobra.Command{
		Use:   "tx",
		Short: "tx runner pattern",
		Long:  `tx_runner`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start to generate tx runner")

			txrunner.Generate(commonPath)
		},
	}

	cmd.PersistentFlags().StringVarP(&commonPath, "common_path", "c", "", "path to common pkg")
	cmd.MarkPersistentFlagRequired("common_path")
	return cmd
}
