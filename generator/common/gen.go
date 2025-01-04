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
	PackageName      string
	IncludeTelemetry bool
	IncludeMetrics   bool
	IncludeLogger    bool
}

func Generate(log, tel, met bool) {
	filePath := os.Getenv("GOFILE")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to parse Go file: %w\n", err))
	}

	packageName := node.Name.Name

	data := MessageData{
		PackageName:      packageName,
		IncludeTelemetry: tel,
		IncludeMetrics:   met,
		IncludeLogger:    log,
	}

	generateFile("database_gen.go", "database_template.tmpl", data)
	generateDir("sqlite.go", "sqlite_template.tmpl", data)
	// generateFile("postgres_gen.go", "postgres_template.tmpl", data)
	generateDir("cassandra.go", "cassandra_template.tmpl", data)
	generateDir("clickhouse.go", "clickhouse_template.tmpl", data)
	generateDir("postgresv2.go", "postgres_templatev2.tmpl", data)

}

func generateDir(fileName, tmplPath string, data MessageData) {

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
