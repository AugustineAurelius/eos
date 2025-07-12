package projectv2

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/AugustineAurelius/eos/pkg/errors"
	myStrings "github.com/AugustineAurelius/eos/pkg/strings"
)

func GenerateDomain(domain Domain) error {
	domainData := domain.ToDomainData()

	for _, entity := range domainData.Entities {
		var uuidImport, timeImport bool
		for _, field := range entity.Fields {
			if field.Type == "uuid" && !uuidImport {
				entity.DoImports = true
				uuidImport = true
				entity.Imports = append(entity.Imports, "github.com/google/uuid")
			}
			if field.Type == "time" && !timeImport {
				entity.DoImports = true
				timeImport = true
				entity.Imports = append(entity.Imports, "time")
			}
		}
		generateDomainFile(fmt.Sprintf("internal/domain/%s.go", strings.ToLower(entity.Name)), "templates/domain.tmpl", entity)
	}
	return nil
}

func generateDomainFile(fileName, tmplPath string, data EntityData) {
	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		errors.FailErr(fmt.Errorf("Failed to read template %s: %v\n", tmplPath, err))
	}

	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"contains": strings.Contains,
		"methodName": func(s string) string {
			if s == "id" {
				return "ID"
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"snakeCase": func(s string) string {
			return myStrings.ToSnakeCase(s)
		},
		"toGoType": toGoType,
	}).Parse(string(tmplContent))
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

func toGoType(s string) string {
	switch s {
	case "uuid":
		return "uuid.UUID"
	case "string":
		return "string"
	case "int":
		return "int"
	case "float":
		return "float"
	case "bool":
		return "bool"
	case "datetime":
		return "time.Time"
	case "date":
		return "time.Time"
	case "time":
		return "time.Time"

	case "enum":
		return "string"
	}
	return s
}

type DomainData struct {
	DomainPackage string
	Imports       []string
	Entities      []EntityData
}

type EntityData struct {
	Name      string
	Fields    []Type
	Imports   []string
	DoImports bool
}

type Type struct {
	Name    string
	Type    string
	IsArray bool
}

func (d *Domain) ToDomainData() DomainData {
	entities := make([]EntityData, 0, len(d.Entities))
	for k, v := range d.Entities {
		entityData := EntityData{
			Name:   k,
			Fields: make([]Type, 0, len(v.Fields)),
		}
		for name, fieldType := range v.Fields {
			entityData.Fields = append(entityData.Fields, Type{Name: name, Type: fieldType.Type, IsArray: fieldType.IsArray})
		}
		entities = append(entities, entityData)
	}

	return DomainData{
		DomainPackage: "domain",
		Entities:      entities,
	}
}
