package generator

import (
	"github.com/spf13/cobra"
)

func init() {

}
func GeneratorCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generator",
		Short: "generator cmd",
		Long:  `generator cmd`,
	}

	cmd.AddCommand(builderCMD(), repositoryCMD(), txCMD(), commonCMD(), projectCMD(), wrapperCMD())
	return cmd
}
