package txrunner

import (
	"embed"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	"github.com/AugustineAurelius/eos/pkg/errors"
	"github.com/AugustineAurelius/eos/pkg/helpers"
)

//go:embed *
var templateFS embed.FS

type MessageData struct {
	PackageName string
	CommonPath  string
}

func Generate(commonPath string) {
	if commonPath == "" {
		errors.FailErr(fmt.Errorf("Missing path to common package\n"))
	}
	if !strings.HasPrefix(commonPath, "/") {
		commonPath = "/" + commonPath
	}
	filePath := os.Getenv("GOFILE")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to parse Go file: %w\n", err))
	}

	packageName := node.Name.Name

	data := MessageData{
		PackageName: packageName,
		CommonPath:  helpers.GetModulePath() + commonPath,
	}

	generateFile("tx_runner.go", "runner_template.tmpl", data)

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
