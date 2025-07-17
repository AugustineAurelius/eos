package repository

import (
	"embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	"github.com/AugustineAurelius/eos/pkg/errors"
	"github.com/AugustineAurelius/eos/pkg/helpers"
	myStrings "github.com/AugustineAurelius/eos/pkg/strings"
)

//go:embed *
var templateFS embed.FS

type Field struct {
	Name   string
	Type   string
	Column string
}

type MessageData struct {
	PackageName    string
	MessageName    string
	TableName      string
	Fields         []Field
	Columns        string
	Placeholders   string
	ModulePath     string
	TxRunnerPath   string
	CommonPath     string
	Imports        []string
	IteratorNumber []string
	WithDefaultID  bool
}

func Generate(structName string, WithDefaultID bool, tableName string) {
	filePath := os.Getenv("GOFILE")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to parse Go file: %w\n", err))
	}
	var imports []string
	for _, imp := range node.Imports {
		imports = append(imports, imp.Path.Value)
	}
	fields, err := parseStruct(node, structName)
	if err != nil {
		errors.FailErr(err)
	}

	packageName := node.Name.Name
	if tableName == "" {
		tableName = structName + "s"
	}

	data := MessageData{
		PackageName:    packageName,
		MessageName:    structName,
		TableName:      tableName,
		Fields:         fields,
		Columns:        strings.Join(getColumns(fields), ", "),
		Placeholders:   strings.Join(getPlaceholders(structName, fields), ", "),
		ModulePath:     helpers.GetPackagePath(),
		Imports:        imports,
		IteratorNumber: []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64"},
		WithDefaultID:  WithDefaultID,
	}

	generateFile("schema.go", "schema_template.tmpl", data)

	generateFile("get_"+strings.ToLower(structName)+"_gen.go", "get_template.tmpl", data)

	generateFile("create_"+strings.ToLower(structName)+"_gen.go", "create_template.tmpl", data)

	generateFile("update_"+strings.ToLower(structName)+"_gen.go", "update_template.tmpl", data)

	generateFile("delete_"+strings.ToLower(structName)+"_gen.go", "delete_template.tmpl", data)

	generateFile("repository_gen.go", "repository_template.tmpl", data)

	generateFile("cursor_gen.go", "cursor_template.tmpl", data)
	// generateFile("helper_gen.go", "helper_template.tmpl", data)
	generateFile("count_gen.go", "count_template.tmpl", data)

}

func generateFile(fileName, tmplPath string, data MessageData) {
	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to read template %s: %v\n", tmplPath, err))
	}

	// Parse the template
	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"contains": strings.Contains,
		"columns_create": func(fields []Field) string {
			cols := make([]string, 0, len(fields))
			for _, field := range fields {
				if data.WithDefaultID && field.Name == "ID" {
					continue
				}
				cols = append(cols, fmt.Sprintf("Column%s%s", data.MessageName, field.Name))
			}
			return strings.Join(cols, ", ")
		},
		"columns": func(fields []Field) string {
			cols := make([]string, 0, len(fields))
			for _, field := range fields {
				cols = append(cols, fmt.Sprintf("Column%s%s", data.MessageName, field.Name))
			}
			return strings.Join(cols, ", ")
		},
		"values_create": func(fields []Field) string {
			cols := make([]string, 0, len(fields))
			for _, field := range fields {
				if data.WithDefaultID && field.Name == "ID" {
					continue
				}
				cols = append(cols, fmt.Sprintf("model.%s", field.Name))
			}
			return strings.Join(cols, ", ")
		},
		"values": func(fields []Field) string {
			cols := make([]string, 0, len(fields))
			for _, field := range fields {
				cols = append(cols, fmt.Sprintf("model.%s", field.Name))
			}
			return strings.Join(cols, ", ")
		},
		"scan": func(fields []Field) string {
			cols := make([]string, 0, len(fields))
			for _, field := range fields {
				cols = append(cols, fmt.Sprintf("&model.%s", field.Name))
			}
			return strings.Join(cols, ", ")
		},
		"Placeholder": func(index int) string { return fmt.Sprintf("$%d", index+1) },
		"snakeCase": func(s string) string {
			return myStrings.ToSnakeCase(s)
		},
	}).Parse(string(tmplContent)) // Convert content to string
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to parse template %s: %v\n", tmplPath, err))
	}

	file, err := os.Create(fileName)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to create file %s: %v\n", fileName, err))
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		errors.FailErr(fmt.Errorf("Failed to execute template %s: %v", tmplPath, err))
	}

	fmt.Printf("Generated %s\n", fileName)
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

func parseStruct(node *ast.File, structName string) ([]Field, error) {
	fields := make([]Field, 0, 8)
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range structType.Fields.List {
				fieldName := field.Names[0].Name
				fieldType := exprToString(field.Type)
				column := fieldName
				fields = append(fields, Field{
					Name:   fieldName,
					Type:   fieldType,
					Column: column,
				})
			}
			break
		}
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("Struct %s not found", structName)
	}

	return fields, nil

}

func getColumns(fields []Field) []string {
	cols := make([]string, 0, len(fields))
	for _, field := range fields {
		cols = append(cols, fmt.Sprintf("Column%s%s", field.Name, field.Name))
	}
	return cols
}

func getPlaceholders(structName string, fields []Field) []string {
	placeholders := make([]string, 0, len(fields))
	for _, field := range fields {
		placeholders = append(placeholders, fmt.Sprintf("%s.%s", strings.ToLower(structName), field.Name))
	}
	return placeholders
}
