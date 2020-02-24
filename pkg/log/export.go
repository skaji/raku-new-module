package log

import (
	base "log"
	"os"
)

var logger Logger = &CoreLogger{
	Level:  3,
	Logger: base.New(os.Stderr, "", base.LstdFlags|base.Llongfile),
}

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

func Close() {
	logger.Close()
}

var debug = os.Getenv("DEBUG") != ""

func Debugf(format string, v ...interface{}) {
	if debug {
		logger.Printf(format, v...)
	}
}
