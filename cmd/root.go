package cmd

import (
	"os"

	"github.com/AugustineAurelius/eos/cmd/compose"
	"github.com/AugustineAurelius/eos/cmd/generator"
	"github.com/AugustineAurelius/eos/cmd/generator/builder"
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

}
