package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/AugustineAurelius/eos/generator/project"
)

func HandleProject() {
	var (
		outputDir   = flag.String("o", "", "output dir path")
		projectName = flag.String("p", "", "name of project")
		url         = flag.String("u", "", "path to repos")
	)

	flag.Parse()

	if *projectName == "" {
		fmt.Println("Error: project name is required (-p flag)")
		flag.Usage()
		os.Exit(1)
	}

	err := project.Generate(project.ProjectData{
		ProjectURL:  *url,
		ProjectName: *projectName,
		Output:      *outputDir,
	})

	if err != nil {
		fmt.Printf("Error generating project: %v\n", err)
		os.Exit(1)
	}
}
