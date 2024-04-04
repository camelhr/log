package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("debug level", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testDebugApp"
		InitGlobalLogger(appName, "debug")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		testMessage := "test message %s %d"
		globalLogger.Debug(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "debug")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})

	t.Run("info level", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testInfoApp"
		InitGlobalLogger(appName, "info")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		testMessage := "test message %s %d"
		globalLogger.Info(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "info")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})

	t.Run("warn level", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testWarnApp"
		InitGlobalLogger(appName, "warn")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		testMessage := "test message %s %d"
		globalLogger.Warn(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "warn")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})

	t.Run("error level", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testErrorApp"
		InitGlobalLogger(appName, "error")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		testMessage := "test message %s %d"
		globalLogger.Error(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "error")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})

	t.Run("panic level", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testPanicApp"
		InitGlobalLogger(appName, "panic")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		testMessage := "test message %s %d"
		// note: this test will cause a panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("the code did not panic")
			}
		}()
		globalLogger.Panic(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "panic")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})

	t.Run("fatal level", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testFatalApp"
		InitGlobalLogger(appName, "fatal")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		testMessage := "test message %s %d"
		// note: This test will cause the program to exit
		if os.Getenv("BE_CRASHER") == "1" {
			globalLogger.Fatal(testMessage, "hello", 7)
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestLogger/fatal level")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err = validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "fatal")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})

	t.Run("additional key-values using 'with' method", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testWithApp"
		InitGlobalLogger(appName, "info")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		newLogger := globalLogger.With("key1", "value1")
		_, ok := newLogger.(*zlWrapper)
		if !ok {
			t.Errorf("With method did not return a *zlWrapper")
		}

		newLogger2 := globalLogger.With("key2", "value2").With("key3", "value3")
		_, ok = newLogger2.(*zlWrapper)
		if !ok {
			t.Errorf("With method did not return a *zlWrapper")
		}

		testMessage := "test message %s %d"
		newLogger2.Info(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "info")
		if err != nil {
			t.Error(err)
		}

		// check that the key values are present in the logged message
		if parsedOutputMap["key2"] != "value2" {
			t.Errorf("Logged message key2 is missing. Got %s, expected value2", parsedOutputMap["key2"])
		}

		if parsedOutputMap["key3"] != "value3" {
			t.Errorf("Logged message key3 is missing. Got %s, expected value3", parsedOutputMap["key3"])
		}

		if len(parsedOutputMap) != 6 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 6", len(parsedOutputMap))
		}
	})

	t.Run("with method does not modify the original logger", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testWithApp"
		InitGlobalLogger(appName, "info")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		newLogger := globalLogger.With("key1", "value1")
		_, ok := newLogger.(*zlWrapper)
		if !ok {
			t.Errorf("With method did not return a *zlWrapper")
		}

		testMessage := "test message %s %d"
		newLogger.Info(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "info")
		if err != nil {
			t.Error(err)
		}

		// check that the key values are present in the logged message
		if parsedOutputMap["key1"] != "value1" {
			t.Errorf("Logged message key1 is missing. Got %s, expected value1", parsedOutputMap["key1"])
		}

		if len(parsedOutputMap) != 5 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 5", len(parsedOutputMap))
		}

		// check that the original logger does not have the key values
		buf.Reset()
		globalLogger.Info(testMessage, "hello", 7)
		parsedOutputMap = make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err = validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "info")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})

	t.Run("with method does not modify the original logger when chaining", func(t *testing.T) {
		initOnce = sync.Once{}
		globalLogger = nil

		// initialize logger
		appName := "testWithApp"
		InitGlobalLogger(appName, "info")

		// create a bytes.Buffer to use as the output destination
		var buf bytes.Buffer

		// set the logger's output to the bytes.Buffer
		globalLogger.SetOutput(&buf)

		newLogger := globalLogger.With("key1", "value1")
		_, ok := newLogger.(*zlWrapper)
		if !ok {
			t.Errorf("With method did not return a *zlWrapper")
		}

		newLogger2 := newLogger.With("key2", "value2")
		_, ok = newLogger2.(*zlWrapper)
		if !ok {
			t.Errorf("With method did not return a *zlWrapper")
		}

		testMessage := "test message %s %d"
		newLogger2.Info(testMessage, "hello", 7)
		expectedMessage := "test message hello 7"

		// parse json from the buffer
		parsedOutputMap := make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err := validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "info")
		if err != nil {
			t.Error(err)
		}

		// check that the key values are present in the logged message
		if parsedOutputMap["key1"] != "value1" {
			t.Errorf("Logged message key1 is missing. Got %s, expected value1", parsedOutputMap["key1"])
		}

		if parsedOutputMap["key2"] != "value2" {
			t.Errorf("Logged message key2 is missing. Got %s, expected value2", parsedOutputMap["key2"])
		}

		if len(parsedOutputMap) != 6 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 6", len(parsedOutputMap))
		}

		// check that the original logger does not have the key values
		buf.Reset()
		globalLogger.Info(testMessage, "hello", 7)
		parsedOutputMap = make(map[string]interface{})
		json.Unmarshal(buf.Bytes(), &parsedOutputMap)

		err = validateLoggedProperties(parsedOutputMap, expectedMessage, appName, "info")
		if err != nil {
			t.Error(err)
		}

		if len(parsedOutputMap) != 4 {
			t.Errorf("Logged message has unexpected number of properties. Got %d, expected 4", len(parsedOutputMap))
		}
	})
}

func validateLoggedProperties(parsedOutputMap map[string]interface{}, expectedMessage string, appName string, level string) error {
	if parsedOutputMap["m"] != expectedMessage {
		return fmt.Errorf("Logged message not found in buffer. Got %s, expected %s", parsedOutputMap["m"], expectedMessage)
	}

	if parsedOutputMap["l"] != level {
		return fmt.Errorf("Logged message level is wrong. Got %s, expected it to be debug", parsedOutputMap["l"])
	}

	if parsedOutputMap["a"] != appName {
		return fmt.Errorf("Logged message app name is wrong. Got %s, expected it to be %s", parsedOutputMap["a"], appName)
	}

	if _, ok := parsedOutputMap["t"]; !ok {
		return fmt.Errorf("Logged message timestamp is missing")
	}

	return nil
}
