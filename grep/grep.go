package grep

import "strings"

// Filter passes strings that match a pattern on the input channel to the output channel.
type Filter struct {
	InputChannel  chan string
	OutputChannel chan string
	match         func(s string) bool
}

// MakeFilter creates a grep Filter.
func MakeFilter(inputChannel, outputChannel chan string, pattern string) Filter {
	filter := Filter{
		InputChannel:  inputChannel,
		OutputChannel: outputChannel,
		match: func(s string) bool {
			return strings.Contains(s, pattern)
		},
	}

	return filter
}

// Start starts a filter's processing.
func (filter *Filter) Start() {
	for inputString := range filter.InputChannel {
		if filter.match(inputString) {
			filter.OutputChannel <- inputString
		}
	}
}
