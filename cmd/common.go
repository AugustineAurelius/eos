package cmd

import (
	"flag"
	"fmt"

	"github.com/AugustineAurelius/eos/generator/common"
)

// HandleCommon handles the common command
func HandleCommon() {
	var (
		loggerEnabled    = flag.Bool("l", false, "add logger to common implementations")
		telemetryEnabled = flag.Bool("t", false, "add telemetry to common implementations")
		metricsEnabled   = flag.Bool("m", false, "add metric to common implementations")
	)

	flag.Parse()

	fmt.Println("start to generate common files")
	common.Generate(*loggerEnabled, *telemetryEnabled, *metricsEnabled)
}
