package generator

import (
	"github.com/AugustineAurelius/eos/cmd/generator/builder"
	"github.com/AugustineAurelius/eos/cmd/generator/repository"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(builder.Cmd)
	Cmd.AddCommand(repository.Cmd)
}

var (
	Cmd = &cobra.Command{
		Use:   "generator",
		Short: "generator cmd",
		Long:  `generator cmd`,
	}
)
