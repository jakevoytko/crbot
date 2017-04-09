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
