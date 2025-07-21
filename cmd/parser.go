package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/AugustineAurelius/eos/generator/parser"
)

func HandleParser() {
	var (
		fileName = flag.String("file", "", "name of the file to parse")
	)

	flag.Parse()

	if *fileName == "" {
		fmt.Println("Error: file name is required (-f flag)")
		flag.Usage()
		os.Exit(1)
	}
	p := parser.NewParser()
	result, err := p.ParseFile(*fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	enums := p.Inspect(result)

	for _, enum := range enums {
		fmt.Printf("enum: %#v\n", enum)
		fmt.Printf("enum.Type.Text: %#v\n", enum.Doc.Text())
		if enum.Type != nil {
			fmt.Printf("type: %#v\n", enum.Type)
		}
	}
}
