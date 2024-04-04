package log

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

func TestDebug(t *testing.T) {
	initOnce = sync.Once{}
	globalLogger = nil
	InitGlobalLogger("test", "debug")

	var buf bytes.Buffer
	SetOutput(&buf)
	Debug("test debug %s", "message")
	if !strings.Contains(buf.String(), "test debug message") {
		t.Errorf("Debug() failed, expected %s, got %s", "test debug message", buf.String())
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	Info("test info %s", "message")
	if !strings.Contains(buf.String(), "test info message") {
		t.Errorf("Info() failed, expected %s, got %s", "test info message", buf.String())
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	Error("test error %s", "message")
	if !strings.Contains(buf.String(), "test error message") {
		t.Errorf("Error() failed, expected %s, got %s", "test error message", buf.String())
	}
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Panic("test panic %s", "message")
}

func TestFatal(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		Fatal("test fatal %s", "message")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestSetOutput(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	Info("test info %s", "message")
	if !strings.Contains(buf.String(), "test info message") {
		t.Errorf("SetOutput() failed, expected %s, got %s", "test info message", buf.String())
	}
}

func TestWith(t *testing.T) {
	logger := With("key", "value")
	if logger == nil {
		t.Errorf("With() failed, got nil")
	}
}
