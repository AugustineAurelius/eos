package common

import (
	"fmt"

	"github.com/AugustineAurelius/eos/generator/common"
	"github.com/spf13/cobra"
)

var TelemetryEnabled, MetricsEnabled bool

var Cmd = &cobra.Command{
	Use:   "common",
	Short: "common files",
	Long:  `common files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to generate common files")

		common.Generate(TelemetryEnabled, MetricsEnabled)
	},
}
