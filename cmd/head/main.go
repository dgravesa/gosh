package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"bitbucket.org/dangravester/gosh/head"
)

func main() {
	// create head filter with command line arguments
	headFilter, fileNames, err := head.NewFilterFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// TODO: allow from stdin if no files specified
	if len(fileNames) < 1 {
		log.Fatal("head: No file specified")
	}

	doneChannel := make(chan struct{})
	defer close(doneChannel)

	// TODO: allow multiple files in succession
	fileName := fileNames[0]
	file, err := os.Open(fileName)
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
			select {
			case inputChannel <- scanner.Text():
			case <-doneChannel:
				return
			}
		}
	}()

	// launch head filter
	headOutputChannel := headFilter.Start(inputChannel, doneChannel)

	// print output channel
	for outputString := range headOutputChannel {
		fmt.Println(outputString)
	}
}
