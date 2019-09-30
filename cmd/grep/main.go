package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bitbucket.org/dangravester/gosh/grep"
	goshio "bitbucket.org/dangravester/gosh/io"
)

var printLineNum = flag.Bool("n", false, "Print line number with output lines")
var invertMatch = flag.Bool("v", false, "Select non-matching lines")

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
	prefixedWriter := goshio.NewPrefixedWriter(os.Stdout)

	fileNames := flag.Args()[1:]

	// we're going to prepend file name to output lines if grepping multiple files
	prefixFileNames := len(fileNames) > 1

	if len(fileNames) == 0 {
		// parse on stdin
		err := filter.Execute(os.Stdin, prefixedWriter)
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
					prefixedWriter.SetPrefix(fileName + ":")
				}

				err = filter.Execute(f, prefixedWriter)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}
