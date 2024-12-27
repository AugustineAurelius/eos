package repositorygen

import (
	"embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

//go:embed repository_gen/*
var templateFS embed.FS

type Field struct {
	Name   string
	Type   string
	Column string
}

type MessageData struct {
	PackageName  string
	MessageName  string
	TableName    string
	Fields       []Field
	Columns      string
	Placeholders string
}

func Generate(structName string) {

	filePath := os.Getenv("GOFILE") // Get the file path from the go:generate directive

	// Parse the Go file to extract the struct
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Failed to parse Go file: %v\n", err)
		return
	}

	// Find the struct definition
	var fields []Field
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
				column := strings.ToLower(fieldName)
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
		fmt.Printf("Struct %s not found in file %s\n", structName, filePath)
		return
	}

	// Generate CRUD methods and constants
	packageName := node.Name.Name
	tableName := strings.ToLower(structName) + "s"

	data := MessageData{
		PackageName:  packageName,
		MessageName:  structName,
		TableName:    tableName,
		Fields:       fields,
		Columns:      strings.Join(getColumns(fields), ", "),
		Placeholders: strings.Join(getPlaceholders(structName, fields), ", "),
	}

	// Generate schema file
	generateFile("schema.go", "templates/schema_template.tmpl", data)
	// Generate CRUD files
	generateFile("get_"+strings.ToLower(structName)+"_gen.go", "repository_gen/get_template.tmpl", data)
	generateFile("create_"+strings.ToLower(structName)+"_gen.go", "repository_gen/create_template.tmpl", data)
	generateFile("update_"+strings.ToLower(structName)+"_gen.go", "repository_gen/update_template.tmpl", data)
	generateFile("delete_"+strings.ToLower(structName)+"_gen.go", "repository_gen/delete_template.tmpl", data)
}

func generateFile(fileName, tmplPath string, data MessageData) {
	// Read the embedded template
	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		fmt.Printf("Failed to read template %s: %v\n", tmplPath, err)
		return
	}

	// Parse the template
	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"columns": func(fields []Field) string {
			var cols []string
			for _, field := range fields {
				cols = append(cols, fmt.Sprintf("Column%s%s", data.MessageName, field.Name))
			}
			return strings.Join(cols, ", ")
		},
		"scanFields": func(fields []Field) string {
			var scanFields []string
			for _, field := range fields {
				scanFields = append(scanFields, fmt.Sprintf("&{{$.MessageName | lower}}.%s", field.Name))
			}
			return strings.Join(scanFields, ", ")
		},
	}).Parse(string(tmplContent)) // Convert content to string
	if err != nil {
		fmt.Printf("Failed to parse template %s: %v\n", tmplPath, err)
		return
	}

	// Create the output file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("Failed to execute template %s: %v\n", tmplPath, err)
		return
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

func getColumns(fields []Field) []string {
	var cols []string
	for _, field := range fields {
		cols = append(cols, fmt.Sprintf("Column%s%s", field.Name, field.Name))
	}
	return cols
}

func getPlaceholders(structName string, fields []Field) []string {
	var placeholders []string
	for _, field := range fields {
		placeholders = append(placeholders, fmt.Sprintf("%s.%s", strings.ToLower(structName), field.Name))
	}
	return placeholders
}
