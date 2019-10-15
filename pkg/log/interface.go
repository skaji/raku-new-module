package log

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
}
