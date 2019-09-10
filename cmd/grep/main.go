package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const patternDescription = "Pattern to match"

// long option
// TODO: pattern should be allowed as first non-option argument
var patternFlag = flag.String("pattern", "", patternDescription)

func init() {
	// short option
	flag.StringVar(patternFlag, "p", "", patternDescription)
}

func main() {
	flag.Parse()

	pattern := *patternFlag
	if len(pattern) == 0 {
		log.Fatal("grep: No pattern specified")
	}

	// for now, handle a single file and return error if one is not specified
	// TODO: make this such that no file specified reads from stdin
	if flag.NArg() != 1 {
		log.Fatal("grep: No file specified")
	}

	fileName := flag.Arg(0)

	// open file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read lines in file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// print line if it contains the pattern
		if strings.Contains(scanner.Text(), pattern) {
			fmt.Println(scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
