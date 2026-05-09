package oasm

import (
	"fmt"
	"time"
)

// Colors for log levels
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

type LoggerType struct {
	name string
}

func Logger(name string) *LoggerType {
	return &LoggerType{
		name: name,
	}
}

func (l *LoggerType) log(color string, message string) {
	currentTime := time.Now()
	// Define the layout for formatting
	layout := "02/01/2006, 3:04:05 pm"

	// Format the current time according to the layout
	currentTimeStr := currentTime.Format(layout)

	fmt.Printf("%s[%s] - [%s] - %s%s\n", color, l.name, currentTimeStr, message, colorReset)
}

// Log logs a message using the LoggerType.
// The message parameter is the message to be logged.
func (l *LoggerType) Log(message string) {
	l.log(colorReset, message)
}

// Error logs an error message using the LoggerType.
// The message parameter is the error message to be logged.
func (l *LoggerType) Error(message string) {
	l.log(colorRed, message)
}

// Success logs a success message using the LoggerType.
func (l *LoggerType) Success(message string) {
	l.log(colorGreen, message)
}

// Warning logs a warning message.
//
// message is the message to be logged.
func (l *LoggerType) Warning(message string) {
	l.log(colorYellow, message)
}

// Debug logs a debug message using the LoggerType.
// The message parameter is the debug message to be logged.
func (l *LoggerType) Debug(message string) {
	l.log(colorPurple, message)
}

// Verbose logs a verbose message using the LoggerType.
// The message parameter is the verbose message to be logged.
func (l *LoggerType) Verbose(message string) {
	l.log(colorCyan, message)
}
