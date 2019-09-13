package grep

import (
	"strings"
)

// Filter passes strings that match a pattern on the input channel to the output channel.
type Filter struct {
	match func(s string) bool
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

// Start turns on a filter's processing.
func (filter *Filter) Start(inputChannel, outputChannel chan string) {
	for inputString := range inputChannel {
		if filter.match(inputString) {
			outputChannel <- inputString
		}
	}
}
