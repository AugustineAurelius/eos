package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AugustineAurelius/eos/pkg/errors"
)

func ValidateFlag(flag string) string {
	if flag == "" {
		errors.FailErr(fmt.Errorf("%s flag is empty\n", flag))
	}
	if !strings.HasPrefix(flag, "/") {
		flag = "/" + flag
	}
	return flag
}

func GetPackagePath() string {
	goModPath, _ := exec.Command("go", "env", "GOMOD").Output()
	goModPathStr := strings.TrimSpace(string(goModPath))

	data, _ := os.ReadFile(goModPathStr)
	moduleLine := strings.Split(string(data), "\n")[0]
	modulePath := strings.TrimPrefix(moduleLine, "module ")

	currentDir, _ := os.Getwd()

	relPath, _ := filepath.Rel(filepath.Dir(goModPathStr), currentDir)

	importPath := filepath.Join(modulePath, relPath)
	return importPath
}

func GetModulePath() string {

	goModPath, _ := exec.Command("go", "env", "GOMOD").Output()
	goModPathStr := strings.TrimSpace(string(goModPath))

	data, _ := os.ReadFile(goModPathStr)
	moduleLine := strings.Split(string(data), "\n")[0]
	modulePath := strings.TrimPrefix(moduleLine, "module ")

	return strings.TrimSpace(modulePath)
}
