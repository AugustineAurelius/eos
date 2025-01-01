package generator

import (
	"github.com/AugustineAurelius/eos/cmd/generator/builder"
	"github.com/AugustineAurelius/eos/cmd/generator/repository"
	txrunner "github.com/AugustineAurelius/eos/cmd/generator/tx_runner"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(builder.Cmd, repository.Cmd, txrunner.Cmd)
}

var (
	Cmd = &cobra.Command{
		Use:   "generator",
		Short: "generator cmd",
		Long:  `generator cmd`,
	}
)
