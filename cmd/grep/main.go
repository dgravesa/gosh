package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"bitbucket.org/dangravester/gosh/grep"
)

func main() {
	// create filter with command line arguments
	grepFilter, fileNames, err := grep.NewFilterFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// for now, handle a single file and return error if one is not specified
	// TODO: make this such that no file specified reads from stdin
	if len(fileNames) < 1 {
		log.Fatal("grep: No file specified")
	}

	// NOTE: unused, but done to adhere to pipeline pattern
	doneChannel := make(chan struct{})
	defer close(doneChannel)

	// open file
	file, err := os.Open(fileNames[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// launch input channel
	// TODO: make this into a utility
	inputChannel := make(chan string)
	go func() {
		defer close(inputChannel)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case inputChannel <- scanner.Text():
			case <-doneChannel:
				return
			}
		}
	}()

	// launch filter
	grepOutputChannel := grepFilter.Start(inputChannel, doneChannel)

	// print outputs of the grep filter
	for outputString := range grepOutputChannel {
		fmt.Println(outputString)
	}
}
