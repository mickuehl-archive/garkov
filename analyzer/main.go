package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mickuehl/garkov"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <model name> <path_to_file> \n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	name := os.Args[1]
	inputFile := os.Args[2]

	// load the model
	model := garkov.New(name, 3)
	defer model.Close()

	// add new training text
	model.Train(inputFile)

	// dump the model
	model.Debug()

	// spill some sentences
	fmt.Println("Markov says: " + model.Sentence(6, 20))
}
