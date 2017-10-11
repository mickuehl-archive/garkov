package main

import (
	"fmt"

	"github.com/mickuehl/garkov"
)

func main() {

	// load the model
	model := garkov.New("test", 3)
	defer model.Close()

	// add new training text
	model.Train("../texts/grim1.txt")
	model.Train("../texts/grim2.txt")
	model.Train("../texts/grim3.txt")

	fmt.Println("\nMarkov says:\n")

	// spill some sentences
	fmt.Println(model.Sentence(6, 20))
	fmt.Println(model.Sentence(6, 20))
	fmt.Println(model.Sentence(6, 20))
	fmt.Println(model.Sentence(6, 20))
	fmt.Println(model.Sentence(6, 20))

	fmt.Println("")
}
