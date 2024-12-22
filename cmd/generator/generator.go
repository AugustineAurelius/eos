package generator

import (
	"github.com/AugustineAurelius/eos/cmd/generator/builder"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(builder.Cmd)
}

var (
	Cmd = &cobra.Command{
		Use:   "generator",
		Short: "generator cmd",
		Long:  `generator cmd`,
	}

	Desination, Source, Package string
)
