package logger

import (
	"log"
	"os"
)

type ILogger interface {
	Info(string)
	Error(string)
}

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
	}
}

func (l *Logger) Info(log string) {
	l.InfoLog.Println(log)
}

func (l *Logger) Error(log string) {
	l.ErrorLog.Println(log)
}
