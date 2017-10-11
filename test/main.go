package main

import (
	"fmt"

	"github.com/mickuehl/garkov"
)

const (
	minWords int = 4
	maxWords int = 40
)

func main() {

	// load the model
	model := garkov.New("test", 2)
	defer model.Close()

	// add new training text
	//model.Train("../texts/grim1.txt") // afther death
	//model.Train("../texts/grim2.txt") // red-riding hood
	//model.Train("../texts/grim3.txt") // Hansel and Gretel

	model.Train("../texts/foo1.nsfw") // random smut

	fmt.Println("\nMarkov says:\n")

	// spill some sentences
	fmt.Println(model.Sentence(minWords, maxWords))
	fmt.Println(model.Sentence(minWords, maxWords))
	fmt.Println(model.Sentence(minWords, maxWords))
	fmt.Println(model.Sentence(minWords, maxWords))
	fmt.Println(model.Sentence(minWords, maxWords))

	fmt.Println("")
}
