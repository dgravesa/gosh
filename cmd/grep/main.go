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

	// open file
	file, err := os.Open(fileNames[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// launch input channel
	inputChannel := make(chan string)
	go func() {
		defer close(inputChannel)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			inputChannel <- scanner.Text()
		}
	}()

	// launch filter
	outputChannel := make(chan string)
	go func() {
		grepFilter.Start(inputChannel, outputChannel)

		// filter returns when input channel is closed
		// since this is the only feeder to the output channel, close it upon completion
		close(outputChannel)
	}()

	// print outputs of the grep filter
	for outputString := range outputChannel {
		fmt.Println(outputString)
	}
}
