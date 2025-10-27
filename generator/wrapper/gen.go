package wrapper

import (
	"embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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

	// Middleware selection flags
	MiddlewareTemplates map[string]bool
	Logging             bool
	Tracing             bool
	NewRelic            bool
	Timeout             bool
	OtelMetrics         bool
	Prometheus          bool
	Retry               bool
	CircuitBreaker      bool
	ContextLogging      bool

	IncludePrivateMethods bool
}

type method struct {
	Name          string
	Signature     string
	Params        []string
	Results       []string
	InputObjects  []obj
	OutputObjects []obj
	HasContext    bool
	HasError      bool
	ErrorParam    string
}

type obj struct {
	Name    string
	Type    string
	ZapType string
}

func Generate(data StructData) error {
	filePath := os.Getenv("GOFILE")
	pkgDir := filepath.Dir(filePath)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse package directory: %w", err)
	}

	var targetPkg *ast.Package
	for pkgName, pkg := range pkgs {
		if !strings.HasSuffix(pkgName, "_test") {
			targetPkg = pkg
			break
		}
	}
	if targetPkg == nil {
		return fmt.Errorf("no non-test package found in directory %s", pkgDir)
	}

	data.PackageName = targetPkg.Name

	for _, file := range targetPkg.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}

			if fn.Name.Name == "" {
				continue
			}

			if !data.IncludePrivateMethods && !isPublicMethod(fn.Name.Name) {
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

			if typeName != data.Name {
				continue
			}

			inputs := formatFieldList(fn.Type.Params)
			outputs := formatFieldList(fn.Type.Results)
			signature := fmt.Sprintf("(%s) (%s)", strings.Join(inputs, ","), strings.Join(outputs, ","))

			errorParam, hasError := hasErrorResult(fn.Type.Results)
			data.Methods = append(data.Methods, method{
				Name:          fn.Name.Name,
				Signature:     signature,
				Params:        formatParamsFieldList(fn.Type.Params),
				Results:       formatResultsFieldList(fn.Type.Results),
				InputObjects:  formatToObjcets(fn.Type.Params),
				OutputObjects: formatToObjcets(fn.Type.Results),
				HasContext:    hasContextParam(fn.Type.Params),
				HasError:      hasError,
				ErrorParam:    errorParam,
			})
		}
	}

	// Initialize middleware templates map if not set
	if data.MiddlewareTemplates == nil {
		data.MiddlewareTemplates = make(map[string]bool)
	}

	// Check if any middleware is selected
	hasSelection := false
	for _, enabled := range data.MiddlewareTemplates {
		if enabled {
			hasSelection = true
			break
		}
	}

	// If no selection, enable all middleware by default
	if !hasSelection {
		data.MiddlewareTemplates["logging"] = true
		data.MiddlewareTemplates["tracing"] = true
		data.MiddlewareTemplates["newrelic"] = true
		data.MiddlewareTemplates["timeout"] = true
		data.MiddlewareTemplates["otel_metrics"] = true
		data.MiddlewareTemplates["prometheus"] = true
		data.MiddlewareTemplates["retry"] = true
		data.MiddlewareTemplates["circuit_breaker"] = true
		data.MiddlewareTemplates["context_logging"] = true
	}

	// Map of middleware template files
	templateFiles := map[string]string{
		"logging":         "templates/logging.tmpl",
		"context_logging": "templates/context_logging.tmpl",
		"tracing":         "templates/tracing.tmpl",
		"newrelic":        "templates/newrelic.tmpl",
		"timeout":         "templates/timeout.tmpl",
		"otel_metrics":    "templates/otel_metrics.tmpl",
		"prometheus":      "templates/prometheus.tmpl",
		"retry":           "templates/retry.tmpl",
		"circuit_breaker": "templates/circuit_breaker.tmpl",
	}

	// Generate base file first
	if err := generateFile("wrapper_"+strings.ToLower(data.Name)+"_gen.go", "templates/base.tmpl", data); err != nil {
		return err
	}

	// Generate middleware files based on selection
	for templateKey, templateFile := range templateFiles {
		if enabled, exists := data.MiddlewareTemplates[templateKey]; exists && enabled {
			if err := appendToFile("wrapper_"+strings.ToLower(data.Name)+"_gen.go", templateFile, data); err != nil {
				return err
			}
		}
	}

	return nil
}

func formatFieldList(fl *ast.FieldList) []string {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}

	var parts []string
	for _, field := range fl.List {
		typeStr := exprToString(field.Type)
		if len(field.Names) == 0 {
			paramName := generateSnakeCaseName(typeStr)
			parts = append(parts, fmt.Sprintf("%s %s", paramName, typeStr))
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
	for _, field := range fl.List {
		typeStr := exprToString(field.Type)
		if len(field.Names) == 0 {
			paramName := generateSnakeCaseName(typeStr)
			parts = append(parts, paramName)
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

func formatToObjcets(fl *ast.FieldList) []obj {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}

	var parts []obj
	for _, field := range fl.List {
		typeStr := exprToString(field.Type)
		if len(field.Names) == 0 {
			paramName := generateSnakeCaseName(typeStr)
			parts = append(parts, obj{
				Name:    paramName,
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

func hasErrorResult(fl *ast.FieldList) (string, bool) {
	if fl == nil || len(fl.List) == 0 {
		return "", false
	}

	for _, field := range fl.List {
		typeStr := exprToString(field.Type)
		if strings.Contains(typeStr, "error") {
			return generateSnakeCaseName(typeStr), true
		}
	}

	return "", false
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
		"zeroValue": func(typeStr string) string {
			switch typeStr {
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"float32", "float64":
				return "0"
			case "bool":
				return "false"
			case "string":
				return "\"\""
			case "error":
				return "nil"
			default:
				if strings.HasPrefix(typeStr, "*") || strings.HasPrefix(typeStr, "[]") || strings.HasPrefix(typeStr, "map[") {
					return "nil"
				}
				return typeStr + "{}"
			}
		},
		"add": func(a, b int) int { return a + b },
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

func appendToFile(fileName, tmplPath string, data StructData) error {
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
		"zeroValue": func(typeStr string) string {
			switch typeStr {
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"float32", "float64":
				return "0"
			case "bool":
				return "false"
			case "string":
				return "\"\""
			case "error":
				return "nil"
			default:
				if strings.HasPrefix(typeStr, "*") || strings.HasPrefix(typeStr, "[]") || strings.HasPrefix(typeStr, "map[") {
					return "nil"
				}
				return typeStr + "{}"
			}
		},
		"add": func(a, b int) int { return a + b },
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", tmplPath, err)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %v", tmplPath, err)
	}

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

// generateSnakeCaseName generates a snake_case name from a type
func generateSnakeCaseName(typeStr string) string {
	// Remove common prefixes and suffixes
	typeStr = strings.TrimPrefix(typeStr, "*")
	typeStr = strings.TrimPrefix(typeStr, "[]")

	// Handle common types
	switch typeStr {
	case "int":
		return "integer"
	case "int8":
		return "integer8"
	case "int16":
		return "integer16"
	case "int32":
		return "integer32"
	case "int64":
		return "integer64"
	case "uint":
		return "unsigned_integer"
	case "uint8":
		return "unsigned_integer8"
	case "uint16":
		return "unsigned_integer16"
	case "uint32":
		return "unsigned_integer32"
	case "uint64":
		return "unsigned_integer64"
	case "float32":
		return "float32"
	case "float64":
		return "float64"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	case "error":
		return "err"
	case "context.Context":
		return "context"
	default:
		// Handle imported types (containing dots)
		if strings.Contains(typeStr, ".") {
			// For imported types like decimal.Decimal, replace dots with underscores
			// and convert to camelCase
			typeStr = strings.ReplaceAll(typeStr, ".", "_")
			var result strings.Builder
			parts := strings.Split(typeStr, "_")
			for i, part := range parts {
				if i == 0 {
					// First part in lowercase
					result.WriteString(strings.ToLower(part))
				} else {
					// Subsequent parts capitalized
					if len(part) > 0 {
						result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
					}
				}
			}
			return result.String()
		}

		// Convert CamelCase to snake_case for local types
		var result strings.Builder
		for i, r := range typeStr {
			if i > 0 && unicode.IsUpper(r) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		}
		return result.String()
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

func isPublicMethod(name string) bool {
	if name == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}
