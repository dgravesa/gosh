package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dgravesa/gosh/grep"
	goshio "github.com/dgravesa/gosh/io"
)

var printLineNumFlag = flag.Bool("n", false, "Print line number with output lines")
var invertMatchFlag = flag.Bool("v", false, "Select non-matching lines")
var recurseDirectoriesFlag = flag.Bool("r", false, "Parse directories recursively")
var ignoreCaseFlag = flag.Bool("i", false, "Ignore case distinctions")

func main() {
	flag.Parse()
	printLineNum := *printLineNumFlag
	invertMatch := *invertMatchFlag
	recurseDirectories := *recurseDirectoriesFlag
	ignoreCase := *ignoreCaseFlag

	if flag.NArg() < 1 {
		log.Fatal("grep: No pattern specified")
	}

	// create grep filter
	pattern := flag.Arg(0)
	filterParams := grep.NewDefaultFilterParams(pattern)
	filterParams.PrintLineNum = printLineNum
	filterParams.InvertMatch = invertMatch
	filterParams.IgnoreCase = ignoreCase
	filter := grep.NewFilter(filterParams)

	// create output writer
	prefixedWriter := goshio.NewPrefixedWriter(os.Stdout)

	inputPaths := flag.Args()[1:]

	// we're going to prepend file name to output lines if grepping multiple files
	prefixFileNames := len(inputPaths) > 1 || recurseDirectories

	if len(inputPaths) == 0 {
		// parse on stdin
		err := filter.Execute(os.Stdin, prefixedWriter)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		var fileNames []string

		if recurseDirectories {
			fileNames = listRegularFilesWithin(inputPaths)
		} else {
			fileNames = checkRegularFiles(inputPaths)
		}

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
