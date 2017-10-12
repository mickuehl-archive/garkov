package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	numParams int = 2
)

func main() {

	if len(os.Args) < (numParams + 1) {
		fmt.Printf("usage: %s <chunk size> <path_to_file> ... \n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	// the file chunck size in bytes
	chunkSize, _ := strconv.Atoi(os.Args[1])
	chunkSize = chunkSize * 1024

	// the file to split
	fileToSplit := os.Args[2]

	// prepare filename parts
	fileDir, fileName := filepath.Split(fileToSplit)
	fileNameBase := strings.Split(fileName, ".")[0]

	fmt.Println(fileDir)
	fmt.Println(fileName)

	// open and read the file
	file, err := os.Open(fileToSplit)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	part := 100 // start with 100 to get a proper lexical order when scanning the directory
	partSize := 0

	// check if the parts folder exists and create it if not
	partsPath := fileDir + fileNameBase

	if _, err = os.Stat(partsPath); os.IsNotExist(err) {
		os.Mkdir(partsPath, os.ModePerm)
	}

	// open the first part file
	partFilePath := buildPartFilePath(fileDir, fileNameBase, part)
	fmt.Println("Wiriting to file " + partFilePath)
	partFile, err := os.Create(partFilePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// read a single line from the input file
		line := scanner.Text()

		// write it to the part file
		partFile.Write([]byte(line + "\n"))

		// update stats
		partSize = partSize + len(line)
		if partSize > chunkSize {
			// close current part file and open the next one
			partFile.Close()

			part = part + 1
			partFilePath := buildPartFilePath(fileDir, fileNameBase, part)
			fmt.Println("Wiriting to file " + partFilePath)

			// new file
			partFile, err = os.Create(partFilePath)
			if err != nil {
				log.Fatal(err)
			}

			// reset the counter
			partSize = 0
		}
	}

}

func buildPartFilePath(dir, baseName string, part int) string {
	return fmt.Sprintf("%v%v/part_%v.txt", dir, baseName, part)
}
