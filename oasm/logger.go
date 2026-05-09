package oasm

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
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
	fmt.Printf("%s[%s] - [%s] - %s%s\n",
		color,
		l.name,
		time.Now().Format("02/01/2006, 3:04:05 pm"),
		fmt.Sprintf(msg, args...),
		colorReset,
	)
}

func (l *LoggerType) Info(msg string, args ...any)    { l.log(slog.LevelInfo, colorReset, msg, args...) }
func (l *LoggerType) Success(msg string, args ...any) { l.log(slog.LevelInfo, colorGreen, msg, args...) }
func (l *LoggerType) Error(msg string, args ...any)   { l.log(slog.LevelError, colorRed, msg, args...) }
func (l *LoggerType) Warning(msg string, args ...any) { l.log(slog.LevelWarn, colorYellow, msg, args...) }
func (l *LoggerType) Debug(msg string, args ...any)   { l.log(slog.LevelDebug, colorPurple, msg, args...) }
func (l *LoggerType) Verbose(msg string, args ...any) { l.log(slog.LevelDebug, colorCyan, msg, args...) }
