package builder

import (
	"fmt"
	"os"

	"github.com/AugustineAurelius/eos/builder"
	"github.com/AugustineAurelius/eos/pkg/errors"
	"github.com/spf13/cobra"
)

var Desination, Source, Package, StructName string

var Cmd = &cobra.Command{
	Use:   "builder",
	Short: "builder pattern",
	Long:  `builder pattern for struct`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to generate builder pattern for struct")

		if StructName == "" {
			fmt.Fprintf(os.Stderr, "Error: %v\n", "use --struct flag with name of your struct")
			os.Exit(1)
		}

		errors.FailErr(builder.Generate(Source, StructName))
	},
}
