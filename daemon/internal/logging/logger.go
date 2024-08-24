package logging

type Logger interface {
	Info(msg string, kvp ...interface{})
	Error(msg string, err error, kvp ...interface{})
	Debug(msg string, kvp ...interface{})
	Close() error
}
