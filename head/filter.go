package head

import (
	"bufio"
	"io"
)

// Filter allows NumLines to pass through when Execute is called
type Filter struct {
	NumLines int
}

// Execute runs the head filter on r, and returns after the number of lines has been read.
func (f *Filter) Execute(r io.Reader, w io.Writer) error {
	numLinesRead := 0

	s := bufio.NewScanner(r)
	for s.Scan() {
		if _, err := w.Write([]byte(s.Text() + "\n")); err != nil {
			return err
		}
		if numLinesRead++; numLinesRead >= f.NumLines {
			break
		}
	}

	return nil
}
