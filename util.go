package main

import "fmt"

// fatal handles a non-recoverable error.
func fatal(msg string, err error) {
	panic(msg + ": " + err.Error())
}

// info prints error information to stdout.
func info(msg string, err error) {
	fmt.Printf(msg+": %v\n", err.Error())
}

// CollapseWhitespace returns an array where splitContent entries are removed,
// starting at the given index. Will return the input if no collapsing is
// necessary, otherwise will return a new slice with the given indices cut out.
func CollapseWhitespace(splitContent []string, startIndex int) []string {
	if startIndex >= len(splitContent) || len(splitContent[startIndex]) > 0 {
		return splitContent
	}

	endIndex := startIndex + 1

	for endIndex < len(splitContent) && len(splitContent[endIndex]) == 0 {
		endIndex++
	}

	return append(splitContent[:startIndex], splitContent[endIndex:]...)
}
