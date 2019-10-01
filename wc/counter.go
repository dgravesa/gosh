package wc

import (
	"bufio"
	"io"
	"strings"
)

// Counter is a struct to accumulate the counts of words, lines, and characters processed.
type Counter struct {
	words int
	lines int
	chars int
}

// NewCounter returns a zero-initialized word/line/char counter.
func NewCounter() *Counter {
	return &Counter{
		words: 0,
		lines: 0,
		chars: 0,
	}
}

// Words returns the number of words processed by a c.
func (c Counter) Words() int {
	return c.words
}

// Lines returns the number of lines processed by a c.
func (c Counter) Lines() int {
	return c.lines
}

// Chars returns the number of characters processed by a c.
func (c Counter) Chars() int {
	return c.chars
}

// Accumulate increments the counts in c with input on r.
func (c *Counter) Accumulate(r io.Reader) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		text := s.Text()
		c.chars += len(text)
		c.words += len(strings.Fields(text))
		c.lines++
	}
}

// Add returns a Counter struct as a sum of the two operands.
func (c Counter) Add(other Counter) Counter {
	return Counter{
		words: c.words + other.words,
		lines: c.lines + other.lines,
		chars: c.chars + other.chars,
	}
}

// Reset zeroes all counts in a Counter.
func (c *Counter) Reset() {
	c.words = 0
	c.lines = 0
	c.chars = 0
}
