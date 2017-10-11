package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mickuehl/garkov"
)

const (
	numParams int = 2
	minWords  int = 4
	maxWords  int = 40
)

func main() {

	if len(os.Args) < (numParams + 1) {
		fmt.Printf("usage: %s <prefix length> <path_to_file>, <path_to_file> ... \n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	// the prefix length
	prefix, _ := strconv.Atoi(os.Args[1])

	// initiate the model
	model := garkov.New("test", prefix)

	// load the files
	i := 2
	for i < len(os.Args) {
		fmt.Println("Reading file: " + os.Args[i])
		model.Train(os.Args[i])
		i = i + 1
	}

	// dump the model
	model.Debug()

}
