package log

import "fmt"

// fatal handles a non-recoverable error.
func Fatal(msg string, err error) {
	panic(msg + ": " + err.Error())
}

// info prints error information to stdout.
func Info(msg string, err error) {
	fmt.Printf(msg+": %v\n", err.Error())
}
