package builder

import (
	. "github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
)

var Desination, Source, Package string

var Cmd = &cobra.Command{
	Use:   "builder",
	Short: "builder pattern",
	Long:  `builder pattern for struct`,
	Run: func(cmd *cobra.Command, args []string) {

		gen(Desination, Package)
	},
}

func gen(dest, pkg string) {
	file := NewFilePathName(dest, pkg)

	file.Comment("Path represent a path with given default values if present")

}
