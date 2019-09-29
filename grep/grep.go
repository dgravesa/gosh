package grep

import (
	"fmt"
	"strings"
)

// Filter passes strings that match a pattern on the input channel to the output channel.
type Filter struct {
	match        func(s string) bool
	printLineNum bool
}

// NewFilter creates a grep filter
func NewFilter(pattern string, printLineNum bool) *Filter {
	return &Filter{
		match: func(s string) bool {
			return strings.Contains(s, pattern)
		},
		printLineNum: printLineNum,
	}
}

// Start turns on a filter's processing.
func (filter *Filter) Start(inputChannel <-chan string) <-chan string {
	outputChannel := make(chan string)

	// launch filter processing goroutine
	go func() {
		// filter is responsible for closing its output channel once complete
		defer close(outputChannel)
		lineNum := 1

		// process until input channel is closed
		for inputString := range inputChannel {
			if filter.match(inputString) {
				// create output string
				var outputString string
				switch filter.printLineNum {
				case true:
					outputString = fmt.Sprintf("%d:%s", lineNum, inputString)
				default:
					outputString = inputString
				}

				// write string to output channel
				outputChannel <- outputString
			}
			lineNum++
		}
	}()

	return outputChannel
}
