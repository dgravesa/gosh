package grep

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// FilterParams is a struct containing parameters to create a grep filter.
type FilterParams struct {
	Pattern      string
	PrintLineNum bool
	InvertMatch  bool
	IgnoreCase   bool
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
		IgnoreCase:   false,
	}
}

// NewFilter creates a grep filter from params.
func NewFilter(params *FilterParams) *Filter {
	invertMatch := params.InvertMatch

	// get pattern to match on
	var pattern string
	if params.IgnoreCase {
		pattern = strings.ToLower(params.Pattern)
	} else {
		pattern = params.Pattern
	}

	return &Filter{
		match: func(input string) bool {
			var s string
			if params.IgnoreCase {
				s = strings.ToLower(input)
			} else {
				s = input
			}
			return strings.Contains(s, pattern) == !invertMatch
		},
		printLineNum: params.PrintLineNum,
	}
}

// Execute runs input from r through the filter and writes to w.
func (filter *Filter) Execute(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	lineNum := 1

	for scanner.Scan() {
		thisLine := scanner.Text() + "\n"

		if filter.match(thisLine) {
			// add line number to string
			if filter.printLineNum {
				thisLine = fmt.Sprintf("%d:%s", lineNum, thisLine)
			}
			// write string
			if _, err := w.Write([]byte(thisLine)); err != nil {
				return err
			}
		}
		lineNum++
	}

	return nil
}
