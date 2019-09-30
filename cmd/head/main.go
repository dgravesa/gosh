package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/dangravester/gosh/head"
)

var numLines = flag.Int("n", 10, "Number of lines to print")

func main() {
	flag.Parse()

	// create head filter
	filter := head.Filter{
		NumLines: *numLines,
	}

	fileNames := flag.Args()

	// print file names before executing if multiple files used as input
	printFileHeader := len(fileNames) > 1

	if len(fileNames) == 0 {
		// run head filter on stdin
		err := filter.Execute(os.Stdin, os.Stdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		for _, fileName := range fileNames {
			if f, err := os.Open(fileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				defer f.Close()
				if printFileHeader {
					fmt.Printf("==> %s <==\n", fileName)
				}
				// run head filter on file
				err = filter.Execute(f, os.Stdout)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}
