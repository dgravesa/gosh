package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"bitbucket.org/dangravester/gosh/grep"
)

var printLineNum = flag.Bool("n", false, "Print line number with output lines")

// TODO: make this a common utility
// newInputChannel returns a readable channel which passes lines from r as they are read.
// Function closes the returned channel once all lines have been read from r.
func newInputChannel(r io.Reader) <-chan string {
	inputChannel := make(chan string)

	go func() {
		defer close(inputChannel)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			inputChannel <- scanner.Text()
		}
	}()

	return inputChannel
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("grep: No pattern specified")
	}

	// create grep filter
	pattern := flag.Arg(0)
	filterParams := grep.NewDefaultFilterParams(pattern)
	filterParams.PrintLineNum = *printLineNum
	filter := grep.NewFilter(filterParams)

	// create output channel
	outputChannel := make(chan string)
	outputDoneChannel := make(chan struct{})

	// launch output listener
	go func() {
		for outputString := range outputChannel {
			fmt.Println(outputString)
		}
		outputDoneChannel <- struct{}{}
	}()

	fileNames := flag.Args()[1:]

	// prepend file name to output strings if grepping multiple files
	prefixFileNames := len(fileNames) > 1

	if len(fileNames) == 0 {
		// use stdin if no file names given
		inputChannel := newInputChannel(os.Stdin)
		filterOutputChannel := filter.Start(inputChannel)

		// forward to output channel
		for outputLine := range filterOutputChannel {
			outputChannel <- outputLine
		}
	} else {
		// parse each file name given
		for _, fileName := range fileNames {
			f, err := os.Open(fileName)

			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: No such file or directory\n", fileName)
			} else {
				inputChannel := newInputChannel(f)
				filterOutputChannel := filter.Start(inputChannel)

				// forward to output channel
				for outputLine := range filterOutputChannel {
					if prefixFileNames {
						outputLine = fileName + ":" + outputLine
					}
					outputChannel <- outputLine
				}
			}
		}
	}

	close(outputChannel)

	// ensure output writer has completed
	<-outputDoneChannel
}
