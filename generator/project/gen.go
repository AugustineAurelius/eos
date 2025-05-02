package project

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed templates/*
var templateFS embed.FS

type ProjectData struct {
	Github      string
	ProjectName string
	Output      string

	fullPath string
}

func Generate(data ProjectData) error {

	data.fullPath = filepath.Join(data.Output, data.ProjectName)
	createDirectories(data.fullPath)

	if err := generateFile("go.mod", "templates/go_mod.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("main.go", "templates/main.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("Makefile", "templates/make.tmpl", data); err != nil {
		return err
	}
	//api
	if err := generateFile("api/api.yaml", "templates/api/openapi.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("api/config.yaml", "templates/api/openapi_cfg.tmpl", data); err != nil {
		return err
	}

	//cmd
	if err := generateFile("cmd/root.go", "templates/cmd/root.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("cmd/serve.go", "templates/cmd/serve.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("cmd/migrate.go", "templates/cmd/migrate.tmpl", data); err != nil {
		return err
	}
	//config
	if err := generateFile("config/manager.go", "templates/config/manager.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("config/postgres.go", "templates/config/postgres.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("config/config.go", "templates/config/config.tmpl", data); err != nil {
		return err
	}
	//pkg
	if err := generateFile("pkg/common/gen.go", "templates/pkg/common.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("pkg/migration/check_version.go", "templates/pkg/migration.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("pkg/logger/zap.go", "templates/pkg/logger.tmpl", data); err != nil {
		return err
	}
	if err := generateFile("pkg/middleware/middlewares.go", "templates/pkg/middleware.tmpl", data); err != nil {
		return err
	}
	//db
	version := time.Now().UTC().Format("20060102150405")
	filename := fmt.Sprintf("%v_%v", version, "init")
	if err := generateFile(fmt.Sprintf("db/postgres/migrations/%s.go", filename), "templates/db/init.tmpl", data); err != nil {
		return err
	}

	//server
	if err := generateFile("server/handler.go", "templates/server/handler.tmpl", data); err != nil {
		return err
	}
	return nil
}

func createDirectories(fullPath string) {
	dirs := []string{
		"api",
		"cmd",
		"config",
		"server",
		"pkg/migration",
		"pkg/logger",
		"pkg/common",
		"pkg/middleware",
		"db/postgres/migrations",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(fullPath, dir), 0755)
		if err != nil {
			panic(err)
		}
	}
}

func generateFile(fileName, tmplPath string, data ProjectData) error {
	fileFullPath := filepath.Join(data.fullPath, fileName)
	if fileExists(fileFullPath) {
		fmt.Printf("file %s alredy exists \n", fileName)
		return nil
	}

	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %v", tmplPath, err)
	}

	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"contains": strings.Contains,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", tmplPath, err)
	}

	file, err := os.Create(fileFullPath)
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
