package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(message string)
	Error(message string, err error)
	Fatal(message string, err error)
}

type AppLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func NewLogger() *AppLogger {
	return &AppLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *AppLogger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *AppLogger) Error(message string, err error) {
	l.errorLogger.Printf("%s: %v", message, err)
}

func (l *AppLogger) Fatal(message string, err error) {
	l.fatalLogger.Fatalf("%s: %v", message, err)
}
