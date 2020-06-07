package log

import (
	"fmt"
	base "log"
	"os"
)

type CoreLogger struct {
	Level  int
	Logger *base.Logger
}

var debugEnabled = os.Getenv("DEBUG") == "1"

func (l *CoreLogger) Fatal(v ...interface{}) {
	l.Logger.Output(l.Level, fmt.Sprint(v...))
}

func (l *CoreLogger) Printf(format string, v ...interface{}) {
	l.Logger.Output(l.Level, fmt.Sprintf(format, v...))
}

func (l *CoreLogger) Print(v ...interface{}) {
	l.Logger.Output(l.Level, fmt.Sprintln(v...))
}

func (l *CoreLogger) Debug(v ...interface{}) {
	if !debugEnabled {
		return
	}
	l.Logger.Output(l.Level, fmt.Sprintln(v...))
}

func (l *CoreLogger) Debugf(format string, v ...interface{}) {
	if !debugEnabled {
		return
	}
	l.Logger.Output(l.Level, fmt.Sprintf(format, v...))
}

func (l *CoreLogger) Close() {
}
