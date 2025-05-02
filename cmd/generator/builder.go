package generator

import (
	"fmt"

	"github.com/AugustineAurelius/eos/generator/builder"
	"github.com/spf13/cobra"
)

func builderCMD() *cobra.Command {
	var desination, source, structName string

	cmd := &cobra.Command{
		Use:   "builder",
		Short: "builder pattern",
		Long:  `builder pattern for struct`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("start to generate builder pattern for struct")
			return builder.Generate(source, structName, desination)
		},
	}

	cmd.PersistentFlags().StringVarP(&desination, "destination", "d", "./gen/generated.go", "place where generate files would placed")
	cmd.PersistentFlags().StringVarP(&source, "source", "s", "./", "read from")
	cmd.PersistentFlags().StringVarP(&structName, "struct", "n", "", "package name")

	cmd.MarkPersistentFlagRequired("struct")
	return cmd

}
