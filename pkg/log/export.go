package log

import (
	base "log"
	"os"
)

var logger Logger = base.New(os.Stderr, "", base.LstdFlags|base.Llongfile)

func Set(l Logger) {
	logger = l
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Println(v ...interface{}) {
	logger.Println(v...)
}
