package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dgravesa/gosh/wc"
)

var printNumWords = flag.Bool("w", false, "Print number of words")
var printNumLines = flag.Bool("l", false, "Print number of lines")
var printNumChars = flag.Bool("c", false, "Print number of characters")

func counterToString(c wc.Counter, lc, wc, cc bool) string {
	arr := []int{}

	appendIf := func(arr []int, val int, c bool) []int {
		if c {
			return append(arr, val)
		}
		return arr
	}

	arr = appendIf(arr, c.Lines(), lc)
	arr = appendIf(arr, c.Words(), wc)
	arr = appendIf(arr, c.Chars(), cc)

	return strings.Trim(fmt.Sprintf("%v", arr), "[]")
}

func main() {
	flag.Parse()

	printSpecified := *printNumWords || *printNumLines || *printNumChars

	// if nothing specified, print them all
	if !printSpecified {
		*printNumWords = true
		*printNumLines = true
		*printNumChars = true
	}

	counts := wc.NewCounter()

	if flag.NArg() == 0 {
		// use stdin
		counts.Accumulate(os.Stdin)
		counterString := counterToString(*counts, *printNumLines, *printNumWords, *printNumChars)
		fmt.Println(counterString)

	} else if flag.NArg() == 1 {
		// use single input file
		fileName := flag.Arg(0)

		if f, err := os.Open(fileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			defer f.Close()
			counts.Accumulate(f)
			counterString := counterToString(*counts, *printNumLines, *printNumWords, *printNumChars)
			fmt.Println(counterString)
		}

	} else {
		// use multiple input files and maintain accumulating total
		fileNames := flag.Args()
		total := wc.NewCounter()

		for _, fileName := range fileNames {
			counts.Reset()

			if f, err := os.Open(fileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				defer f.Close()
				counts.Accumulate(f)
				counterString := counterToString(*counts, *printNumLines, *printNumWords, *printNumChars)
				fmt.Printf("%s %s\n", counterString, fileName)
			}

			*total = total.Add(*counts)
		}

		// print totals
		totalString := counterToString(*total, *printNumLines, *printNumWords, *printNumChars)
		fmt.Printf("%s total\n", totalString)
	}
}
