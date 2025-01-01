package txrunner

import (
	"fmt"

	txrunner "github.com/AugustineAurelius/eos/generator/tx_runner"
	"github.com/spf13/cobra"
)

var CommonPath string

var Cmd = &cobra.Command{
	Use:   "tx",
	Short: "tx runner pattern",
	Long:  `tx_runner`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to generate tx runner")

		txrunner.Generate(CommonPath)
	},
}
