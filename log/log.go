package log

import "fmt"

// Fatal handles a non-recoverable error.
func Fatal(msg string, err error) {
	panic(msg + ": " + err.Error())
}

// Info prints error information to stdout.
func Info(msg string, err error) {
	fmt.Printf(msg+": %v\n", err.Error())
}
