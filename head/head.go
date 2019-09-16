package head

import "flag"

const numLinesDescription string = "Number of lines to allow"

// Filter passes a limited number of strings through once started.
type Filter struct {
	numLines int
}

// NewFilter creates a new head filter to pass N lines through.
func NewFilter(n int) *Filter {
	return &Filter{numLines: n}
}

// NewFilterFromArgs creates a new head filter from a list of arguments
func NewFilterFromArgs(arguments []string) (*Filter, []string, error) {

	// create flag parser
	headFlags := flag.NewFlagSet("head", flag.ContinueOnError)

	// create options
	numLines := headFlags.Int("lines", 10, numLinesDescription)
	headFlags.IntVar(numLines, "n", 10, numLinesDescription)

	// parse arguments
	if err := headFlags.Parse(arguments); err != nil {
		return nil, headFlags.Args(), err
	}

	return NewFilter(*numLines), headFlags.Args(), nil
}

// NumLines returns the number of lines that filter allows to pass through.
func (filter *Filter) NumLines() int {
	return filter.numLines
}

// Start begins processing of the head filter, which stops after the number of lines has been processed.
func (filter *Filter) Start(inputChannel <-chan string, done <-chan struct{}) <-chan string {
	outputChannel := make(chan string)

	go func() {
		defer close(outputChannel)
		numLinesPassed := 0

		for inputString := range inputChannel {
			// block until send or done signal received
			select {
			case outputChannel <- inputString:
			case <-done:
				return
			}

			numLinesPassed++

			if numLinesPassed == filter.numLines {
				break
			}
		}
	}()

	return outputChannel
}
