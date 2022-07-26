package logger

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Fatal(string)
	Error(string)
	Warning(string)
	Info(string)
	Debug(string)
	Trace(string)
}

type FileLogger struct {
	writer io.Writer
}

func New(filePath string) FileLogger {
	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		panic(err)
	}

	return FileLogger{
		writer: writer,
	}
}

func (l *FileLogger) Fatal(msg string) {
	io.WriteString(l.writer, fmt.Sprintf("Fatal: %s\n", msg))
}

func (l *FileLogger) Error(msg string) {
	io.WriteString(l.writer, fmt.Sprintf("Error: %s\n", msg))
}

func (l *FileLogger) Warning(msg string) {
	io.WriteString(l.writer, fmt.Sprintf("Warning: %s\n", msg))
}

func (l *FileLogger) Info(msg string) {
	io.WriteString(l.writer, fmt.Sprintf("Info: %s\n", msg))
}

func (l *FileLogger) Debug(msg string) {
	io.WriteString(l.writer, fmt.Sprintf("Debug: %s\n", msg))
}

func (l *FileLogger) Trace(msg string) {
	io.WriteString(l.writer, fmt.Sprintf("Trace: %s\n", msg))
}
