package logger

const (
	Trace   LogLevel = "trace"
	Debug   LogLevel = "debug"
	Info    LogLevel = "info"
	Warning LogLevel = "warning"
	Error   LogLevel = "error"
	Fatal   LogLevel = "fatal"
)

type Config struct {
	Level LogLevel `env:"LEVEL" envDefault:"info"`
}

type LogLevel string
