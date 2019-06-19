package log

import (
	base "log"
	"os"
)

var logger Logger = base.New(os.Stderr, "", base.LstdFlags|base.Llongfile)

// Set is
func Set(l Logger) {
	logger = l
}

// Fatal is
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

// Printf is
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// Println is
func Println(v ...interface{}) {
	logger.Println(v...)
}
