package grep

import (
	"fmt"
	"strings"
)

// FilterParams is a struct containing parameters to create a grep filter.
type FilterParams struct {
	Pattern      string
	PrintLineNum bool
	InvertMatch  bool
}

// Filter passes strings that match a pattern on the input channel to the output channel.
type Filter struct {
	match        func(s string) bool
	printLineNum bool
}

// NewDefaultFilterParams returns a default set of filter parameters.
func NewDefaultFilterParams(pattern string) *FilterParams {
	return &FilterParams{
		Pattern:      pattern,
		PrintLineNum: false,
		InvertMatch:  false,
	}
}

// NewFilter creates a grep filter from params
func NewFilter(params *FilterParams) *Filter {
	invertMatch := params.InvertMatch

	return &Filter{
		match: func(s string) bool {
			return strings.Contains(s, params.Pattern) == !invertMatch
		},
		printLineNum: params.PrintLineNum,
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
