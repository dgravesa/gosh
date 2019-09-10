package main

import (
	"flag"
	"log"
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
}
