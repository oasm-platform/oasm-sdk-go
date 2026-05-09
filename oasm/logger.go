package oasm

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	// colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

type LoggerType struct {
	handler slog.Handler
	name    string
}

func NewLogger(name string) *LoggerType {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	return &LoggerType{
		handler: slog.NewTextHandler(os.Stdout, opts),
		name:    name,
	}
}

func (l *LoggerType) log(level slog.Level, color string, msg string, args ...any) {
	currentTimeStr := time.Now().Format("2006-01-02 15:04:05")
	formattedMsg := fmt.Sprintf(msg, args...)

	fmt.Printf("%s[%s] %s | %-5s | %s%s\n",
		color, l.name, currentTimeStr, level.String(), formattedMsg, colorReset)
}

func (l *LoggerType) ErrorE(msg string, err error, args ...any) {
	l.log(slog.LevelError, colorRed, msg+": %v", append(args, err)...)
}

func (l *LoggerType) Info(msg string, args ...any) { l.log(slog.LevelInfo, colorReset, msg, args...) }
func (l *LoggerType) Success(msg string, args ...any) {
	l.log(slog.LevelInfo, colorGreen, msg, args...)
}
func (l *LoggerType) Error(msg string, args ...any) { l.log(slog.LevelError, colorRed, msg, args...) }
func (l *LoggerType) Warning(msg string, args ...any) {
	l.log(slog.LevelWarn, colorYellow, msg, args...)
}
func (l *LoggerType) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, colorPurple, msg, args...)
}
func (l *LoggerType) Verbose(msg string, args ...any) {
	l.log(slog.LevelDebug, colorCyan, msg, args...)
}
