package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

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
	generateFile("schema.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/schema_template.tmpl", data)
	// Generate CRUD files
	generateFile("get_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/get_template.tmpl", data)
	generateFile("create_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/create_template.tmpl", data)
	generateFile("update_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/update_template.tmpl", data)
	generateFile("delete_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/delete_template.tmpl", data)
}

func generateFile(fileName, tmplURL string, data MessageData) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", tmplURL, nil)
	if err != nil {
		fmt.Printf("Failed to create request for %s: %v\n", tmplURL, err)
		return
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to fetch template %s: %v\n", tmplURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch template %s: status code %d\n", tmplURL, resp.StatusCode)
		return
	}

	// Read the response body into a string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body for %s: %v\n", tmplURL, err)
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
	}).Parse(string(body)) // Convert body to string
	if err != nil {
		fmt.Printf("Failed to parse template %s: %v\n", tmplURL, err)
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
		fmt.Printf("Failed to execute template %s: %v\n", tmplURL, err)
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
