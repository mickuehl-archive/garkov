package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mickuehl/garkov"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <path_to_file> \n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	// open and read the file
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read or create the dictionary
	dict := garkov.OpenDictionary("test")
	defer dict.Close()

	fmt.Println(dict)

	// read the file and analyze it
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dict.AddSentence(scanner.Text())
	}

	//fmt.Println(dict)

	// some cleanup?
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
