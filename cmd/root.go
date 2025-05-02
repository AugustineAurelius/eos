package cmd

import (
	"os"

	"github.com/AugustineAurelius/eos/cmd/generator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eos",
	Short: "framework for generation different stuff",
	Long: `In ancient Greek mythology and religion, Eos is the goddess and personification of the dawn,
	who rose each morning from her home at the edge of the river Oceanus to deliver light and disperse the night.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generator.GeneratorCMD())
}
