package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	infoMsg     = "an info message"
	warnMsg     = "a warn message"
	errorMsg    = "an error message"
	debugMsg    = "a debug message"
	newErrorMsg = "new error"
	printfMsg   = "my printf message"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, envLogLevel, "LOG_LEVEL")
	assert.EqualValues(t, envLogOutput, "LOG_OUTPUT")
}

func TestMemSinkClose(t *testing.T) {
	sink = &MemorySink{new(bytes.Buffer)}
	result := sink.Close()
	assert.Nil(t, result)
}

func TestMemSinkSync(t *testing.T) {
	sink = &MemorySink{new(bytes.Buffer)}
	result := sink.Sync()
	assert.Nil(t, result)
}

func TestGetOutputReturnsOutput(t *testing.T) {
	t.Setenv("LOG_OUTPUT", "logoutputtest")
	result := getOutput()
	assert.NotNil(t, result)
	assert.EqualValues(t, "logoutputtest", result)
}

func TestGetOutputReturnsStdoutByDefault(t *testing.T) {
	t.Setenv("LOG_OUTPUT", " ")
	result := getOutput()
	assert.EqualValues(t, "stdout", result)
}

func TestGetLevelReturnsInfoForInvalidLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "invalid")
	result := getLevel()
	assert.EqualValues(t, "info", result.String())
}

func TestGetLoggerReturnsLogger(t *testing.T) {
	myLogger := GetLogger()
	assert.NotNil(t, myLogger)
}

func extractLog() map[string]any {
	output := sink.String()
	m := make(map[string]any)
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	return m
}

func TestInfoWritesInfo(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	initLogger(true, "")
	Info(infoMsg)
	m := extractLog()
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], infoMsg)
}

func TestInfoWithFieldWritesInfoWithFields(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	initLogger(true, "")
	Info(infoMsg, Field{
		Key:   "id",
		Value: "123",
	})
	m := extractLog()
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], infoMsg)
	assert.EqualValues(t, m["id"], "123")
}

func TestErrorWritesError(t *testing.T) {
	t.Setenv("LOG_LEVEL", "error")
	initLogger(true, "")
	Error(errorMsg, errors.New(newErrorMsg))
	m := extractLog()
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], errorMsg)
	assert.EqualValues(t, m["error"], newErrorMsg)
}

func TestErrorWithFieldWritesErrorWithField(t *testing.T) {
	t.Setenv("LOG_LEVEL", "error")
	initLogger(true, "")
	Error(errorMsg, errors.New(newErrorMsg), Field{
		Key:   "id",
		Value: "123",
	})
	m := extractLog()
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], errorMsg)
	assert.EqualValues(t, m["error"], newErrorMsg)
	assert.EqualValues(t, m["id"], "123")
}

func TestDebugWritesDebug(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	initLogger(true, "")
	Debug(debugMsg)
	m := extractLog()
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], debugMsg)
}

func TestDebugWithFieldWritesDebugWithField(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	initLogger(true, "")
	Debug(debugMsg, Field{
		Key:   "id",
		Value: "123",
	})
	m := extractLog()
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], debugMsg)
	assert.EqualValues(t, m["id"], "123")
}

func TestPrintPrints(t *testing.T) {
	initLogger(true, "")
	log.Print("a", "b")
	m := extractLog()
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "ab")
}

func TestPrintfPrints(t *testing.T) {
	initLogger(true, "")
	log.Printf(printfMsg)
	m := extractLog()
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], printfMsg)
}

func TestPrintfWithFormatPrints(t *testing.T) {
	initLogger(true, "")
	log.Printf("my %s message", "formatted")
	m := extractLog()
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my formatted message")
}

func TestWarnWritesWarn(t *testing.T) {
	initLogger(true, "")
	Warn(warnMsg)
	m := extractLog()
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], warnMsg)
}

func TestWarnWithFieldWritesWarnWithFields(t *testing.T) {
	initLogger(true, "")
	Warn(warnMsg, Field{
		Key:   "id",
		Value: "123",
	})
	m := extractLog()
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], warnMsg)
	assert.EqualValues(t, m["id"], "123")
}

func TestWriteInfoWritesInfo(t *testing.T) {
	initLogger(true, "")
	written, writeErr := log.Write([]byte(infoMsg))
	m := extractLog()
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(infoMsg), written)
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], infoMsg)
}

func TestWriteWarnWritesWarn(t *testing.T) {
	initLogger(true, "")
	written, writeErr := log.Write([]byte(warnMsg))
	m := extractLog()
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(warnMsg), written)
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], warnMsg)
}

func TestWriteErrorWritesError(t *testing.T) {
	initLogger(true, "")
	written, writeErr := log.Write([]byte(errorMsg))
	m := extractLog()
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(errorMsg), written)
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], errorMsg)
}

func TestWriteDebugWritesDebug(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	initLogger(true, "")
	written, writeErr := log.Write([]byte(debugMsg))
	m := extractLog()
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(debugMsg), written)
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], debugMsg)
}

func TestAddtoLogListDoesNotOverflow(t *testing.T) {
	ClearLogList()
	for range 1100 {
		addToLogList("Info", "I was here")
	}
	assert.EqualValues(t, logListMaxLength, len(GetLogList()))
}

func TestAddtoLogListRetainsEntries(t *testing.T) {
	ClearLogList()
	addToLogList("Info", "One")
	addToLogList("Warn", "Two")
	addToLogList("Error", "Three")
	l := GetLogList()
	assert.EqualValues(t, "Info", l[0].LogLevel)
	assert.EqualValues(t, "One", l[0].LogMessage)
	assert.EqualValues(t, "Warn", l[1].LogLevel)
	assert.EqualValues(t, "Two", l[1].LogMessage)
	assert.EqualValues(t, "Error", l[2].LogLevel)
	assert.EqualValues(t, "Three", l[2].LogMessage)
}

func TestClearLogListClearsLogList(t *testing.T) {
	ClearLogList()
	addToLogList("Info", "One")
	addToLogList("Warn", "Two")
	addToLogList("Error", "Three")
	l1 := GetLogList()
	ClearLogList()
	l2 := GetLogList()
	assert.EqualValues(t, 3, len(l1))
	assert.EqualValues(t, 0, len(l2))
}

func TestGetLogListReturnsCopy(t *testing.T) {
	ClearLogList()
	addToLogList("Info", "One")
	l := GetLogList()
	l[0].LogMessage = "changed"

	result := GetLogList()
	assert.EqualValues(t, "One", result[0].LogMessage)
}

func TestAddToLogListHandlesConcurrentWrites(t *testing.T) {
	ClearLogList()
	var wg sync.WaitGroup
	for range 1100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addToLogList("Info", "I was here")
		}()
	}
	wg.Wait()

	assert.EqualValues(t, logListMaxLength, len(GetLogList()))
}

func TestDebugAddsToLogList(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	initLogger(true, "")
	ClearLogList()

	Debug(debugMsg)

	l := GetLogList()
	assert.EqualValues(t, 1, len(l))
	assert.EqualValues(t, "Debug", l[0].LogLevel)
	assert.EqualValues(t, debugMsg, l[0].LogMessage)
}

func TestDebugfWritesDebugWithFormat(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	initLogger(true, "")
	Debugf("my debug message: %v", "A")
	m := extractLog()
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my debug message: A")
}

func TestInfofWritesInfoWithFormat(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	initLogger(true, "")
	Infof("my info message: %v", "A")
	m := extractLog()
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my info message: A")
}

func TestWarnfWritesWarnWithFormat(t *testing.T) {
	t.Setenv("LOG_LEVEL", "warn")
	initLogger(true, "")
	Warnf("my warn message: %v", "A")
	m := extractLog()
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my warn message: A")
}

func TestErrorfWritesErrorWithFormat(t *testing.T) {
	t.Setenv("LOG_LEVEL", "error")
	initLogger(true, "")
	Errorf("my error message: %v", "A")
	m := extractLog()
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my error message: A")
}

func TestErrorWithNilErrorWritesErrorWithoutErrorField(t *testing.T) {
	t.Setenv("LOG_LEVEL", "error")
	initLogger(true, "")
	Error(errorMsg, nil)
	m := extractLog()
	assert.EqualValues(t, m["level"], "error")
	assert.EqualValues(t, m["msg"], errorMsg)
	assert.NotContains(t, m, "error")
}

func TestLogToFile(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "app.log")
	initLogger(true, logFile)
	t.Cleanup(func() {
		initLogger(true, "")
	})
	Infof("my log message: %v", "A")

	data, err := os.ReadFile(logFile)
	assert.Nil(t, err)
	assert.Contains(t, string(data), "\"level\":\"info\"")
	assert.Contains(t, string(data), "\"msg\":\"my log message: A\"")
}
