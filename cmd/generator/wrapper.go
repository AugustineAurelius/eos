package generator

import (
	"github.com/AugustineAurelius/eos/generator/wrapper"
	"github.com/spf13/cobra"
)

func wrapperCMD() *cobra.Command {

	var name string
	cmd := &cobra.Command{
		Use:   "wrapper",
		Short: "wrapper generator",

		RunE: func(cmd *cobra.Command, args []string) error {

			return wrapper.Generate(wrapper.StructData{
				Name: name,
			})
		},
	}

	cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "name of the struct")

	cmd.MarkPersistentFlagRequired("name")

	return cmd
}
