package generator

import (
	"fmt"

	"github.com/AugustineAurelius/eos/generator/common"
	"github.com/spf13/cobra"
)

func commonCMD() *cobra.Command {
	var loggerEnabled, telemetryEnabled, metricsEnabled bool

	cmd := &cobra.Command{
		Use:   "common",
		Short: "common files",
		Long:  `common files`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start to generate common files")
			common.Generate(loggerEnabled, telemetryEnabled, metricsEnabled)
		},
	}

	cmd.PersistentFlags().BoolVarP(&telemetryEnabled, "telemetry", "t", false, "add telemetry to common implementations")
	cmd.PersistentFlags().BoolVarP(&metricsEnabled, "metrics", "m", false, "add metric to common implementations")
	cmd.PersistentFlags().BoolVarP(&loggerEnabled, "logger", "l", false, "add logger to common implementations")

	return cmd
}
