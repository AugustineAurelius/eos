package repositorygen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"text/template"
)

type Field struct {
	Name       string
	Type       string
	Column     string
	ForeignKey string
}

type MessageData struct {
	PackageName   string
	MessageName   string
	TableName     string
	Fields        []Field
	JoinFields    []Field
	Columns       string
	Placeholders  string
	JoinTable     string
	JoinCondition string
	JoinColumns   string
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
				fk := ""
				if field.Tag != nil {
					tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
					fk = tag.Get("fk") // Extract foreign key from the `fk` tag
				}
				fields = append(fields, Field{
					Name:       fieldName,
					Type:       fieldType,
					Column:     column,
					ForeignKey: fk,
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

	var joinFields []Field
	var joinTable, joinCondition, joinColumns string

	for _, field := range fields {
		if field.ForeignKey != "" {
			// Handle foreign key
			refTable, refColumn := parseForeignKey(field.ForeignKey)
			joinTable = strings.ToLower(refTable) + "s"
			joinCondition = fmt.Sprintf("%s.%s = %s.%s", tableName, field.Column, joinTable, strings.ToLower(refColumn))

			// Add fields from the referenced table
			for _, m := range node.Decls {
				genDecl, ok := m.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || typeSpec.Name.Name != refTable {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					for _, f := range structType.Fields.List {
						if f.Names[0].Name != refColumn { // Exclude the foreign key column itself
							joinFields = append(joinFields, Field{
								Name:   refTable + f.Names[0].Name,
								Type:   exprToString(f.Type),
								Column: joinTable + "." + strings.ToLower(f.Names[0].Name),
							})
							joinColumns += joinTable + "." + strings.ToLower(f.Names[0].Name) + ", "
						}
					}
					break
				}
			}
		}
	}

	if joinColumns != "" {
		joinColumns = strings.TrimSuffix(joinColumns, ", ")
	}

	data := MessageData{
		PackageName:   packageName,
		MessageName:   structName,
		TableName:     tableName,
		Fields:        fields,
		JoinFields:    joinFields,
		Columns:       strings.Join(getColumns(fields), ", "),
		Placeholders:  strings.Join(getPlaceholders(structName, fields), ", "),
		JoinTable:     joinTable,
		JoinCondition: joinCondition,
		JoinColumns:   joinColumns,
	}

	generateFile("schema.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/schema_template.tmpl", data)
	// Generate CRUD files
	generateFile("get_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/get_template.tmpl", data)
	generateFile("create_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/create_template.tmpl", data)
	generateFile("update_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/update_template.tmpl", data)
	generateFile("delete_"+strings.ToLower(structName)+"_gen.go", "https://raw.githubusercontent.com/AugustineAurelius/eos/repository_gen/delete_template.tmpl", data)
}

func generateFile(fileName, tmplFile string, data MessageData) {
	tmpl, err := template.New(tmplFile).Funcs(template.FuncMap{
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
	}).ParseFiles(tmplFile)
	if err != nil {
		fmt.Printf("Failed to parse template %s: %v\n", tmplFile, err)
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	if err := tmpl.ExecuteTemplate(file, tmplFile, data); err != nil {
		fmt.Printf("Failed to execute template %s: %v\n", tmplFile, err)
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

func parseForeignKey(foreignKey string) (string, string) {
	parts := strings.Split(foreignKey, ".")
	if len(parts) != 2 {
		panic("invalid foreign key format")
	}
	return parts[0], parts[1]
}
