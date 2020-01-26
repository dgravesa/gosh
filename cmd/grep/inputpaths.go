package main

import (
	"fmt"
	"os"
	"path"
)

func checkRegularFiles(paths []string) []string {
	var regpaths []string

	for _, p := range paths {
		pathstat, err := os.Stat(p)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
		} else {
			switch mode := pathstat.Mode(); {
			case !mode.IsRegular():
				fmt.Fprintf(os.Stderr, "%s: not a regular file\n", p)
			default:
				regpaths = append(regpaths, p)
			}
		}
	}

	return regpaths
}

func listRegularFilesWithin(paths []string) []string {
	var regpaths []string

	for _, p := range paths {
		if pathstat, err := os.Stat(p); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
		} else {
			switch mode := pathstat.Mode(); {

			// add path to list
			case mode.IsRegular():
				regpaths = append(regpaths, p)

			// recursively list regular files within
			case mode.IsDir():
				if dir, err := os.Open(p); err != nil {
					fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
				} else if subpaths, err := dir.Readdirnames(0); err != nil {
					fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
				} else {
					fullsubpaths := prefixPaths(p, subpaths)
					regsubpaths := listRegularFilesWithin(fullsubpaths)
					regpaths = append(regpaths, regsubpaths...)
				}
			}
		}
	}

	return regpaths
}

func prefixPaths(prefix string, paths []string) []string {
	var prefixedpaths []string
	for _, p := range paths {
		prefixedpath := path.Join(prefix, p)
		prefixedpaths = append(prefixedpaths, prefixedpath)
	}
	return prefixedpaths
}
