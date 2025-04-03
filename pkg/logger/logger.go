package logger

import (
	"log"
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *Logger) Info(format string, v ...any) {
	l.logger.Printf(ColorBlue+"[INFO] "+ColorReset+format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	l.logger.Printf(ColorYellow+"[WARN] "+ColorReset+format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.logger.Printf(ColorRed+"[ERROR] "+ColorReset+format, v...)
}

func (l *Logger) Success(format string, v ...any) {
	l.logger.Printf(ColorGreen+"[SUCCESS] "+ColorReset+format, v...)
}

func (l *Logger) Debug(format string, v ...any) {
	l.logger.Printf(ColorCyan+"[DEBUG] "+ColorReset+format, v...)
}

func (l *Logger) Fatal(format string, v ...any) {
	l.logger.Fatalf(ColorRed+"[FATAL] "+ColorReset+format, v...)
}
