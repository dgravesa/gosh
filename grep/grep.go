package grep

import (
	"flag"
	"fmt"
	"strings"
)

const (
	patternDescription      string = "Pattern to match"
	printLineNumDescription string = "Print line number with output lines"
)

// Filter passes strings that match a pattern on the input channel to the output channel.
type Filter struct {
	match        func(s string) bool
	printLineNum bool
}

// NewFilter creates a basic pattern-matching grep Filter.
func NewFilter(pattern string) *Filter {
	filter := Filter{
		match: func(s string) bool {
			return strings.Contains(s, pattern)
		},
	}

	return &filter
}

// NewFilterFromArgs creates a grep Filter from an arguments list.
// On success, the resulting Filter and any unused arguments are returned, and error is set to nil.
// On failure, an error is returned.
func NewFilterFromArgs(arguments []string) (*Filter, []string, error) {

	// create flag parser
	grepFlags := flag.NewFlagSet("grep", flag.ContinueOnError)

	// create long options
	pattern := grepFlags.String("pattern", "", patternDescription)
	printLineNum := grepFlags.Bool("line-number", false, printLineNumDescription)

	// create short options
	grepFlags.StringVar(pattern, "p", "", patternDescription)
	grepFlags.BoolVar(printLineNum, "n", false, printLineNumDescription)

	// parse arguments
	if err := grepFlags.Parse(arguments); err != nil {
		return nil, grepFlags.Args(), err
	}

	// check for pattern specified
	if *pattern == "" {
		return nil, grepFlags.Args(), fmt.Errorf("grep: No pattern specified")
	}

	// create filter
	filter := Filter{
		match: func(s string) bool {
			return strings.Contains(s, *pattern)
		},
		printLineNum: *printLineNum,
	}

	return &filter, grepFlags.Args(), nil
}

// Start turns on a filter's processing.
func (filter *Filter) Start(inputChannel <-chan string, done <-chan struct{}) <-chan string {
	outputChannel := make(chan string)

	// launch filter processing goroutine
	go func() {
		// filter is responsible for closing its output channel once complete
		defer close(outputChannel)
		lineNum := 1

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

				// block until string sends or done signal is received by closing the done channel
				select {
				case outputChannel <- outputString:
				case <-done:
					return
				}
			} else {
				// check for done signal
				select {
				case <-done:
					return
				default:
				}
			}
			lineNum++
		}
	}()

	return outputChannel
}
