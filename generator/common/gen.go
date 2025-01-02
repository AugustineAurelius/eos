package common

import (
	"embed"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	"github.com/AugustineAurelius/eos/pkg/errors"
)

//go:embed *
var templateFS embed.FS

type MessageData struct {
	PackageName string
}

func Generate() {
	filePath := os.Getenv("GOFILE")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to parse Go file: %w\n", err))
	}

	packageName := node.Name.Name

	data := MessageData{
		PackageName: packageName,
	}

	generateFile("database_gen.go", "database_template.tmpl", data)
	generateFile("sqlite_gen.go", "sqlite_template.tmpl", data)
	generateFile("postgres_gen.go", "postgres_template.tmpl", data)
	generateFile("cassandra_gen.go", "cassandra_template.tmpl", data)
	generateFile("clickhouse_gen.go", "clickhouse_template.tmpl", data)

}

func generateFile(fileName, tmplPath string, data MessageData) {
	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to read template %s: %v\n", tmplPath, err))
	}

	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"lower": strings.ToLower,
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
