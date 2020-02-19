# gosh

## Summary

This repository is a small collection of utilities that mimic their Linux command line counterparts.
I've started this repository both as a way for me to replace functionality that I miss the most when working on Windows command prompt
and as a way to better familiarize myself with Golang.

## List of current commands and features

### grep
* only basic matching is currently supported
* supports filtering on `STDIN`
* `-r` handle directories recursively
* `-n` print line numbers
* `-v` invert match
* `-i` ignore case

### wc
* `-w, -c, -l` all implemented for getting word, character, and line counts selectively
* supports aggregation on `STDIN`

### head
* supports filtering on `STDIN`
* `-n` number of lines to print
