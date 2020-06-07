package log

import (
	base "log"
	"os"
)

var logger Logger = &CoreLogger{
	Level:  3,
	Logger: base.New(os.Stderr, "", base.LstdFlags|base.Lshortfile),
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

func Print(v ...interface{}) {
	logger.Print(v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Close() {
	logger.Close()
}
