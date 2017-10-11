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

	if len(os.Args) <= (numParams + 1) {
		fmt.Printf("usage: %s <model depth> <path_to_file>, <path_to_file> ... \n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	depth, _ := strconv.Atoi(os.Args[1])

	// load the model
	model := garkov.New("test", depth)
	defer model.Close()

	// add new training text
	//model.Train("../texts/grim1.txt") // afther death
	//model.Train("../texts/grim2.txt") // red-riding hood
	//model.Train("../texts/grim3.txt") // Hansel and Gretel

	model.Train("../texts/war_and_piece.txt")
	model.Train("../texts/anna_karenina.txt")

	// dump the model
	//model.Debug()

	fmt.Println("\nMarkov says:\n")

	i := 0
	for i < num {
		fmt.Println(model.Sentence(minWords, maxWords))
		i = i + 1
	}

	fmt.Println("")

}
