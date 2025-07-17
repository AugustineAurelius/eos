package cmd

import (
	"flag"
	"fmt"
	"os"

	repositorygen "github.com/AugustineAurelius/eos/generator/repository"
)

func HandleRepository() {
	var (
		structName    = flag.String("type", "", "name of the struct for which would be generated repo")
		withDefaultID = flag.Bool("default_id", false, "add id to create")
		table         = flag.String("table", "", "name of the table for which would be generated repo")
	)

	flag.Parse()

	if *structName == "" {
		fmt.Println("Error: struct name is required (-t flag)")
		flag.Usage()
		os.Exit(1)
	}

	if *table == "" {
		fmt.Println("Error: table name is required (-t flag)")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("start to generate repository for struct")
	repositorygen.Generate(*structName, *withDefaultID, *table)
}
