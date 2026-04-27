package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	oldPrintf := fmtPrintf
	defer func() { fmtPrintf = oldPrintf }()
	fmtPrintf = func(format string, a ...interface{}) (n int, err error) {
		return fmt.Fprintf(&buf, format, a...)
	}
	f()
	return buf.String()
}

var fmtPrintf = fmt.Printf

func TestLogger_Log(t *testing.T) {
	log := Logger("test")
	output := captureOutput(func() {
		log.Log("test message")
	})

	if !strings.Contains(output, "[test]") {
		t.Errorf("expected output to contain [test], got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "-") {
		t.Errorf("expected output to contain timestamp separator, got: %s", output)
	}
	// Default log should not have color codes
	if strings.Contains(output, "\033[") {
		t.Errorf("default log should not contain color codes, got: %s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	log := Logger("test")
	output := captureOutput(func() {
		log.Error("error message")
	})

	if !strings.Contains(output, "[test]") {
		t.Errorf("expected output to contain [test], got: %s", output)
	}
	if !strings.Contains(output, "error message") {
		t.Errorf("expected output to contain 'error message', got: %s", output)
	}
	if !strings.Contains(output, colorRed) {
		t.Errorf("error log should contain red color code, got: %s", output)
	}
	if !strings.Contains(output, colorReset) {
		t.Errorf("error log should contain reset color code, got: %s", output)
	}
}

func TestLogger_Success(t *testing.T) {
	log := Logger("test")
	output := captureOutput(func() {
		log.Success("success message")
	})

	if !strings.Contains(output, "[test]") {
		t.Errorf("expected output to contain [test], got: %s", output)
	}
	if !strings.Contains(output, "success message") {
		t.Errorf("expected output to contain 'success message', got: %s", output)
	}
	if !strings.Contains(output, colorGreen) {
		t.Errorf("success log should contain green color code, got: %s", output)
	}
	if !strings.Contains(output, colorReset) {
		t.Errorf("success log should contain reset color code, got: %s", output)
	}
}

func TestLogger_Warning(t *testing.T) {
	log := Logger("test")
	output := captureOutput(func() {
		log.Warning("warning message")
	})

	if !strings.Contains(output, "[test]") {
		t.Errorf("expected output to contain [test], got: %s", output)
	}
	if !strings.Contains(output, "warning message") {
		t.Errorf("expected output to contain 'warning message', got: %s", output)
	}
	if !strings.Contains(output, colorYellow) {
		t.Errorf("warning log should contain yellow color code, got: %s", output)
	}
	if !strings.Contains(output, colorReset) {
		t.Errorf("warning log should contain reset color code, got: %s", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	log := Logger("test")
	output := captureOutput(func() {
		log.Debug("debug message")
	})

	if !strings.Contains(output, "[test]") {
		t.Errorf("expected output to contain [test], got: %s", output)
	}
	if !strings.Contains(output, "debug message") {
		t.Errorf("expected output to contain 'debug message', got: %s", output)
	}
	if !strings.Contains(output, colorPurple) {
		t.Errorf("debug log should contain purple color code, got: %s", output)
	}
	if !strings.Contains(output, colorReset) {
		t.Errorf("debug log should contain reset color code, got: %s", output)
	}
}

func TestLogger_Verbose(t *testing.T) {
	log := Logger("test")
	output := captureOutput(func() {
		log.Verbose("verbose message")
	})

	if !strings.Contains(output, "[test]") {
		t.Errorf("expected output to contain [test], got: %s", output)
	}
	if !strings.Contains(output, "verbose message") {
		t.Errorf("expected output to contain 'verbose message', got: %s", output)
	}
	if !strings.Contains(output, colorCyan) {
		t.Errorf("verbose log should contain cyan color code, got: %s", output)
	}
	if !strings.Contains(output, colorReset) {
		t.Errorf("verbose log should contain reset color code, got: %s", output)
	}
}

func TestLogger_EmptyMessage(t *testing.T) {
	log := Logger("test")
	output := captureOutput(func() {
		log.Log("")
	})

	if !strings.Contains(output, "[test]") {
		t.Errorf("expected output to contain [test], got: %s", output)
	}
}

func TestLoggerName(t *testing.T) {
	tests := []struct {
		name     string
		logName  string
		expected string
	}{
		{"custom name", "myLogger", "myLogger"},
		{"another name", "test", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := Logger(tt.logName)
			output := captureOutput(func() {
				log.Log("message")
			})

			if !strings.Contains(output, "["+tt.expected+"]") {
				t.Errorf("expected output to contain [%s], got: %s", tt.expected, output)
			}
		})
	}
}
