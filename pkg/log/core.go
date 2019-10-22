package log

import (
	"fmt"
	base "log"
)

type CoreLogger struct {
	Level  int
	Logger *base.Logger
}

func (l *CoreLogger) Fatal(v ...interface{}) {
	l.Logger.Output(l.Level, fmt.Sprint(v...))
}

func (l *CoreLogger) Printf(format string, v ...interface{}) {
	l.Logger.Output(l.Level, fmt.Sprintf(format, v...))
}

func (l *CoreLogger) Println(v ...interface{}) {
	l.Logger.Output(l.Level, fmt.Sprintln(v...))
}
