package builder

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/AugustineAurelius/eos/pkg/errors"
	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

func Generate(Source, StructName string) error {
	goPackage := os.Getenv("GOPACKAGE")

	f := NewFile(goPackage)
	f.PackageComment("Code generated by generator, DO NOT EDIT.")

	structType := parseStruct(Source, StructName)

	builderStructName := strings.ToLower(StructName) + "builder"

	f.Type().Id(builderStructName).Struct(
		Id("inner").Op("*").Id(StructName),
	)

	f.Func().Id("New" + StructName + "Builder").Params().Op("*").Id(builderStructName).Block(
		Return().Op("&").Id(builderStructName).Block(),
	)

	builderMethod := Id("b").Op("*").Id(builderStructName)

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)

		switch v := field.Type().(type) {
		case *types.Basic:

			f.Func().Params(builderMethod).
				Id("Set" + field.Name()).Params(Id(field.Name()).Id(v.String())).
				Block(
					Id("b").Dot("inner").Dot(field.Name()).Op("=").Id(field.Name()),
				)

		case *types.Named:
			typeName := v.Obj()

			f.Func().Params(builderMethod).
				Id("Set" + field.Name()).Params(Id(field.Name()).Qual(typeName.Pkg().Path(), typeName.Name())).
				Block(
					Id("b").Dot("inner").Dot(field.Name()).Op("=").Id(field.Name()),
				)

		case *types.Slice:
			switch s := v.Elem().(type) {
			case *types.Basic:

				f.Func().Params(builderMethod).
					Id("Set" + field.Name()).Params(Id(field.Name()).Index().Id(s.String())).
					Block(
						Id("b").Dot("inner").Dot(field.Name()).Op("=").Id(field.Name()),
					)

				f.Func().Params(builderMethod).
					Id("AddOneTo" + field.Name()).Params(Id("one").Id(s.String())).
					Block(
						Id("b").Dot("inner").Dot(field.Name()).Op("=").Append(Id("b").Dot("inner").Dot(field.Name()), Id("one")),
					)

			case *types.Named:
				typeName := s.Obj()

				f.Func().Params(builderMethod).
					Id("Set" + field.Name()).Params(Id(field.Name()).Index().Qual(typeName.Pkg().Path(), typeName.Name())).
					Block(
						Id("b").Dot("inner").Dot(field.Name()).Op("=").Id(field.Name()),
					)

				f.Func().Params(builderMethod).
					Id("AddOneTo" + field.Name()).Params(Id("one").Qual(typeName.Pkg().Path(), typeName.Name())).
					Block(
						Id("b").Dot("inner").Dot(field.Name()).Op("=").Append(Id("b").Dot("inner").Dot(field.Name()), Id("one")),
					)

			}

		default:
			return fmt.Errorf("struct field type not hanled: %T", v)
		}
	}

	f.Func().Params(builderMethod).
		Id("Build").Params().Id(StructName).
		Block(
			Return(Op("*").Id("b").Dot("inner")),
		)

	goFile := os.Getenv("GOFILE")
	ext := filepath.Ext(goFile)
	baseFilename := goFile[0 : len(goFile)-len(ext)]
	targetFilename := baseFilename + "_gen.go"

	return f.Save(targetFilename)
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		errors.FailErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func parseStruct(Source, StructName string) *types.Struct {
	pkg := loadPackage(Source)
	obj := pkg.Types.Scope().Lookup(StructName)
	if obj == nil {
		errors.FailErr(fmt.Errorf("%s not found in declared types of %s", StructName, pkg))
	}

	if _, ok := obj.(*types.TypeName); !ok {
		errors.FailErr(fmt.Errorf("%v is not a named type", obj))
	}

	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		errors.FailErr(fmt.Errorf("type %v is not a struct", obj))
	}

	return structType
}