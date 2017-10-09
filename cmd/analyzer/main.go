package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mickuehl/garkov"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <model> <path_to_file> \n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	modelName := os.Args[1]
	inputFile := os.Args[2]

	// load the model
	model := garkov.OpenModel(modelName)
	defer model.Close()

	// add new training text
	model.TrainModel(inputFile)
}
