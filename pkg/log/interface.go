package log

type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Fatal(v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Close()
}
