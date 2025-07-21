package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Parser struct {
	fileSet *token.FileSet
}

func NewParser() *Parser {
	return &Parser{fileSet: token.NewFileSet()}
}

func (p *Parser) ParseFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(p.fileSet, fileName, nil, parser.ParseComments)
}

func (p *Parser) Inspect(f ast.Node) map[string]*ast.TypeSpec {
	result := make(map[string]*ast.TypeSpec)
	ast.Inspect(f, func(n ast.Node) bool {
		switch val := n.(type) {
		case *ast.GenDecl:
			copyGenDeclCommentsToSpecs(val)

		case *ast.Ident:
			if val.Obj == nil {
				return true
			}
			// fmt.Printf("Node: %#v\n", val.Obj)
			if val.Obj.Kind == ast.Typ {
				if ts, ok := val.Obj.Decl.(*ast.TypeSpec); ok {
					// fmt.Printf("Type: %+v\n", ts)
					// isEnum := isTypeSpecEnum(ts)
					// Only store documented enums
					// if isEnum {
					// fmt.Printf("EnumType: %T\n", ts.Type)
					result[val.Name] = ts
					// }
				}
			}

		}

		return true
	})
	return result
}

func copyGenDeclCommentsToSpecs(x *ast.GenDecl) {
	if x.Doc != nil {
		for _, spec := range x.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				if s.Doc == nil {
					s.Doc = x.Doc
				}
			case *ast.ValueSpec:
				if s.Doc == nil {
					s.Doc = x.Doc
				}
			}
		}
	}
}
