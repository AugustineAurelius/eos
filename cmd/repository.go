package cmd

import (
	"flag"
	"fmt"
	"os"

	repositorygen "github.com/AugustineAurelius/eos/generator/repository"
)

func HandleRepository() {
	var (
		structName    = flag.String("t", "", "name of the struct for which would be generated repo")
		withDefaultID = flag.Bool("i", false, "add id to create")
	)

	flag.Parse()

	if *structName == "" {
		fmt.Println("Error: struct name is required (-t flag)")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("start to generate repository for struct")
	repositorygen.Generate(*structName, *withDefaultID)
}
