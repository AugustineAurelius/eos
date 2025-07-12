package projectv2

import (
	"embed"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed templates/*
var templateFS embed.FS

func Generate(pathToSpec string) error {
	spec, err := os.ReadFile(pathToSpec)
	if err != nil {
		return err
	}

	var specData Spec

	err = yaml.Unmarshal(spec, &specData)
	if err != nil {
		return err
	}

	createDirectories("")

	GenerateDomain(specData.Domain)

	return nil
}

func createDirectories(fullPath string) {
	dirs := []string{
		"internal/domain",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(fullPath, dir), 0755)
		if err != nil {
			panic(err)
		}
	}
}
