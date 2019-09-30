package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"bitbucket.org/dangravester/gosh/grep"
)

var printLineNum = flag.Bool("n", false, "Print line number with output lines")
var invertMatch = flag.Bool("v", false, "Select non-matching lines")

type prefixedWriter struct {
	Receiver io.Writer
	Prefix   []byte
}

func (w *prefixedWriter) Write(p []byte) (n int, err error) {
	return w.Receiver.Write([]byte(append(w.Prefix, p...)))
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
	filterParams.InvertMatch = *invertMatch
	filter := grep.NewFilter(filterParams)

	// create output writer
	decorator := &prefixedWriter{
		Receiver: os.Stdout,
	}

	fileNames := flag.Args()[1:]

	// we're going to prepend file name to output lines if grepping multiple files
	prefixFileNames := len(fileNames) > 1

	if len(fileNames) == 0 {
		// parse on stdin
		err := filter.Execute(os.Stdin, decorator)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		// parse each file name given
		for _, fileName := range fileNames {
			if f, err := os.Open(fileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				defer f.Close()
				if prefixFileNames {
					decorator.Prefix = []byte(fileName + ":")
				}

				err = filter.Execute(f, decorator)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}
