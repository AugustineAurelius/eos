package cmd

import (
	"os"

	"github.com/AugustineAurelius/eos/cmd/compose"
	"github.com/AugustineAurelius/eos/cmd/generator"
	"github.com/AugustineAurelius/eos/cmd/generator/builder"
	"github.com/AugustineAurelius/eos/cmd/generator/common"
	"github.com/AugustineAurelius/eos/cmd/generator/repository"
	txrunner "github.com/AugustineAurelius/eos/cmd/generator/tx_runner"
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

	flags()
	rootCmd.AddCommand(compose.Cmd)

	rootCmd.AddCommand(generator.Cmd)

}

func flags() {
	compose.Cmd.Flags().BoolVarP(&compose.Postgres, "postgres", "p", false, "add postgres to compose file")
	compose.Cmd.Flags().BoolVarP(&compose.App, "app", "a", false, "add app Dockerfile and to compose file")
	compose.Cmd.Flags().BoolVarP(&compose.Kafka, "kafka", "k", false, "add kafka, zookeeper and kowl to compose file")

	builder.Cmd.Flags().StringVarP(&builder.Desination, "destination", "d", "./gen/generated.go", "place where generate files would placed")
	builder.Cmd.Flags().StringVarP(&builder.Source, "source", "s", "./", "read from")
	builder.Cmd.Flags().StringVarP(&builder.Package, "packege", "p", "generated", "package name")
	builder.Cmd.Flags().StringVarP(&builder.StructName, "struct", "n", "", "package name")

	repository.Cmd.Flags().StringVarP(&repository.StructName, "type", "t", "", "name of the struct for which would be generated repo")
	repository.Cmd.Flags().StringVarP(&repository.TxRunenrPath, "tx_path", "p", "", "path to txrunner impl")
	repository.Cmd.Flags().StringVarP(&repository.CommonPath, "common_path", "c", "", "path to common")
	repository.Cmd.Flags().BoolVarP(&repository.WithTX, "with_tx", "x", false, "does repo would be with transaction logic")
	repository.Cmd.MarkPersistentFlagRequired("common_path")

	txrunner.Cmd.Flags().StringVarP(&txrunner.CommonPath, "common_path", "c", "", "path to common")
	txrunner.Cmd.MarkPersistentFlagRequired("common_path")

	common.Cmd.Flags().BoolVarP(&common.TelemetryEnabled, "telemetry", "t", false, "add metric to common implementations")
	common.Cmd.Flags().BoolVarP(&common.MetricsEnabled, "metrics", "m", false, "add metric to common implementations")
	common.Cmd.Flags().BoolVarP(&common.LoggerEnabled, "logger", "l", false, "add metric to common implementations")

}
