package wrapper

import (
	"embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

//go:embed *
var templateFS embed.FS

type StructData struct {
	Name        string
	PackageName string
	Methods     []method
}

type method struct {
	Name          string
	Signature     string
	Params        []string
	InputObjects  []obj
	OutputObjects []obj
	HasContext    bool
}

type obj struct {
	Name    string
	Type    string
	ZapType string
}

func Generate(data StructData) error {
	filePath := os.Getenv("GOFILE")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse Go file: %w", err)
	}

	data.PackageName = node.Name.Name

	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}
		r, size := utf8.DecodeRuneInString(fn.Name.Name)
		if r == utf8.RuneError && size <= 1 {
			continue
		}
		if unicode.IsLower(r) {
			continue
		}
		recvType := fn.Recv.List[0].Type
		var typeName string
		switch rt := recvType.(type) {
		case *ast.Ident:
			typeName = rt.Name
		case *ast.StarExpr:
			if ident, ok := rt.X.(*ast.Ident); ok {
				typeName = ident.Name
			}
		}

		inputs := formatFieldList(fn.Type.Params)
		outputs := formatFieldList(fn.Type.Results)

		var signature string

		signature = fmt.Sprintf("(%s) (%s)", strings.Join(inputs, ","), strings.Join(outputs, ","))

		if typeName == data.Name {
			data.Methods = append(data.Methods, method{
				Name:          fn.Name.Name,
				Signature:     signature,
				Params:        formatParamsFieldList(fn.Type.Params),
				InputObjects:  formatToObjcets(fn.Type.Params),
				OutputObjects: formatToObjcets(fn.Type.Results),

				HasContext: hasContextParam(fn.Type.Params),
			})

		}
	}
	if err = generateFile("wrapper_"+strings.ToLower(data.Name)+"_gen.go", "wrap.tmpl", data); err != nil {
		return err
	}
	return nil
}

func formatFieldList(fl *ast.FieldList) []string {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}

	var parts []string
	for i, field := range fl.List {
		typeStr := exprToString(field.Type)
		if len(field.Names) == 0 {
			parts = append(parts, fmt.Sprintf("param%d %s", i, typeStr))
			continue
		}

		names := make([]string, len(field.Names))
		for i, name := range field.Names {
			names[i] = name.Name + " " + typeStr
		}
		parts = append(parts, names...)
	}

	return parts
}

func formatParamsFieldList(fl *ast.FieldList) []string {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}

	var parts []string
	for _, field := range fl.List {
		typeStr := exprToString(field.Type)
		if len(field.Names) == 0 {
			parts = append(parts, typeStr)
			continue
		}

		names := make([]string, len(field.Names))
		for i, name := range field.Names {
			names[i] = name.Name
		}
		parts = append(parts, names...)
	}

	return parts
}

func formatResultsFieldList(fl *ast.FieldList) []string {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}

	var parts []string
	for i, field := range fl.List {
		if len(field.Names) == 0 {
			parts = append(parts, fmt.Sprintf("param%d", i))
			continue
		}

	}

	return parts
}

func formatToObjcets(fl *ast.FieldList) []obj {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}

	var parts []obj
	for i, field := range fl.List {
		typeStr := exprToString(field.Type)
		if len(field.Names) == 0 {
			parts = append(parts, obj{
				Name:    fmt.Sprintf("param%d", i),
				Type:    typeStr,
				ZapType: getZapType(typeStr),
			})
			continue
		}

		names := make([]obj, len(field.Names))
		for i, name := range field.Names {
			names[i] = obj{
				Name:    name.Name,
				Type:    typeStr,
				ZapType: getZapType(typeStr),
			}
		}
		parts = append(parts, names...)
	}

	return parts
}

func hasContextParam(fl *ast.FieldList) bool {
	if fl == nil || len(fl.List) == 0 {
		return false
	}

	for _, field := range fl.List {
		if strings.Contains(exprToString(field.Type), "context.Context") {
			return true
		}

	}

	return false
}

func generateFile(fileName, tmplPath string, data StructData) error {

	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %v", tmplPath, err)
	}

	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"contains": strings.Contains,
		"join":     strings.Join,
		"firstToLower": func(s string) string {
			r, size := utf8.DecodeRuneInString(s)
			if r == utf8.RuneError && size <= 1 {
				return s
			}
			lc := unicode.ToLower(r)
			if r == lc {
				return s
			}
			return string(lc) + s[size:]
		},
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", tmplPath, err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", fileName, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %v", tmplPath, err)
	}

	fmt.Printf("Generated %s\n", fileName)
	return nil
}

func exprToString(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return exprToString(v.X) + "." + v.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(v.Elt)
	case *ast.StarExpr:
		return "*" + exprToString(v.X)
	default:
		return "unknown"
	}
}

func getZapType(typeStr string) string {
	switch typeStr {
	case "int":
		return "Int"
	case "int8":
		return "Int8"
	case "int16":
		return "Int16"
	case "int32":
		return "Int32"
	case "int64":
		return "Int64"

	case "uint":
		return "Uint"
	case "uint8":
		return "Uint8"
	case "uint16":
		return "Uint16"
	case "uint32":
		return "Uint32"
	case "uint64":
		return "Uint64"

	case "float32":
		return "Float32"
	case "float64":
		return "Float64"

	case "complex64":
		return "Complex64"
	case "complex128":
		return "Complex128"

	case "string":
		return "String"
	case "bool":
		return "Bool"
	case "error":
		return "Error"

	case "time.Time":
		return "Time"
	case "time.Duration":
		return "Duration"

	case "interface{}":
		return "Any"

	default:
		if strings.HasPrefix(typeStr, "*") {
			return "Any"
		}

		if strings.HasPrefix(typeStr, "[]") || strings.Contains(typeStr, "[") {
			return "Any"
		}

		if strings.HasPrefix(typeStr, "map[") {
			return "Any"
		}

		return "Any"
	}
}
