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
		fileOrDir := os.Args[i]

		fi, err := os.Stat(fileOrDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		if fi.IsDir() {
			fmt.Println("Scanning directory: " + fileOrDir)

			fileList := []string{}
			filepath.Walk(fileOrDir, func(path string, f os.FileInfo, err error) error {
				fileList = append(fileList, path)
				return nil
			})

			fileList = fileList[1:] // remove the first entry, its the directory itself

			for _, file := range fileList {
				fmt.Println("Reading file: " + file)
				model.Build(file)
			}

		} else {
			fmt.Println("Reading file: " + fileOrDir)
			model.Build(fileOrDir)
		}

		i = i + 1
	}

	// dump the model
	model.Debug()
	fmt.Println("")
}
