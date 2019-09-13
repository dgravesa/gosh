package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"bitbucket.org/dangravester/gosh/grep"
)

const patternDescription = "Pattern to match"

// long option
// TODO: pattern should be allowed as first non-option argument
var patternFlag = flag.String("pattern", "", patternDescription)

func init() {
	// short option
	flag.StringVar(patternFlag, "p", "", patternDescription)
}

func main() {
	flag.Parse()

	pattern := *patternFlag
	if len(pattern) == 0 {
		log.Fatal("grep: No pattern specified")
	}

	// for now, handle a single file and return error if one is not specified
	// TODO: make this such that no file specified reads from stdin
	if flag.NArg() != 1 {
		log.Fatal("grep: No file specified")
	}

	fileName := flag.Arg(0)

	// open file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// create the grep filter
	grepFilter := grep.NewFilter(pattern)

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
