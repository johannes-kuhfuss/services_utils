package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	os.Setenv("LOG_OUTPUT", "logoutputtest")
	result := getOutput()
	assert.NotNil(t, result)
	assert.EqualValues(t, "logoutputtest", result)
}

func TestGetLoggerReturnsLogger(t *testing.T) {
	myLogger := GetLogger()
	assert.NotNil(t, myLogger)
}

func TestInfoWritesInfo(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	initLogger(true)
	Info("my info message")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my info message")
}

func TestInfoWithFieldWritesInfoWithFields(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	initLogger(true)
	Info("my info message", Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my info message")
	assert.EqualValues(t, m["id"], "123")
}

func TestErrorWritesError(t *testing.T) {
	os.Setenv("LOG_LEVEL", "error")
	initLogger(true)
	Error("my error message", errors.New("new error"))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my error message")
	assert.EqualValues(t, m["error"], "new error")
}

func TestErrorWithFieldWritesErrorWithField(t *testing.T) {
	os.Setenv("LOG_LEVEL", "error")
	initLogger(true)
	Error("my error message", errors.New("new error"), Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my error message")
	assert.EqualValues(t, m["error"], "new error")
	assert.EqualValues(t, m["id"], "123")
}

func TestDebugWritesDebug(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	initLogger(true)
	Debug("my debug message")
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my debug message")
}

func TestDebugWithFieldWritesDebugWithField(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	initLogger(true)
	Debug("my debug message", Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my debug message")
	assert.EqualValues(t, m["id"], "123")
}

func TestPrintPrints(t *testing.T) {
	initLogger(true)
	log.Print("a", "b")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "[a b]")
}

func TestPrintfPrints(t *testing.T) {
	initLogger(true)
	log.Printf("my printf message")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my printf message")
}

func TestPrintfWithFormatPrints(t *testing.T) {
	initLogger(true)
	log.Printf("my %s message", "formatted")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my formatted message")
}

func TestWarnWritesWarn(t *testing.T) {
	initLogger(true)
	Warn("my warn message")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my warn message")
}

func TestWarnWithFieldWritesWarnWithFields(t *testing.T) {
	initLogger(true)
	Warn("my warn message", Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my warn message")
	assert.EqualValues(t, m["id"], "123")
}

func TestWriteInfoWritesInfo(t *testing.T) {
	initLogger(true)
	logMessage := "the is an info message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is an info message")
}

func TestWriteWarnWritesWarn(t *testing.T) {
	initLogger(true)
	logMessage := "the is a warning message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is a warning message")
}

func TestWriteErrorWritesError(t *testing.T) {
	initLogger(true)
	logMessage := "the is an error message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is an error message")
}

func TestWriteDebugWritesDebug(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	initLogger(true)
	logMessage := "the is a debug message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is a debug message")
}

func TestAddtoLogListDoesNotOverflow(t *testing.T) {
	loglist = nil
	for i := 0; i < 1100; i++ {
		addToLogList("Info", "I was here")
	}
	assert.EqualValues(t, logListMaxLength, len(loglist))
}

func TestAddtoLogListRetainsEntries(t *testing.T) {
	loglist = nil
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
	loglist = nil
	addToLogList("Info", "One")
	addToLogList("Warn", "Two")
	addToLogList("Error", "Three")
	l1 := GetLogList()
	ClearLogList()
	l2 := GetLogList()
	assert.EqualValues(t, 3, len(l1))
	assert.EqualValues(t, 0, len(l2))
}

func TestDebugfWritesDebugWithFormat(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	initLogger(true)
	Debugf("my debug message: %v", "A")
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my debug message: A")
}

func TestInfofWritesInfoWithFormat(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	initLogger(true)
	Infof("my info message: %v", "A")
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my info message: A")
}

func TestWarnfWritesWarnWithFormat(t *testing.T) {
	os.Setenv("LOG_LEVEL", "warn")
	initLogger(true)
	Warnf("my warn message: %v", "A")
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my warn message: A")
}

func TestErrorfWritesErrorWithFormat(t *testing.T) {
	os.Setenv("LOG_LEVEL", "error")
	initLogger(true)
	Errorf("my error message: %v", "A")
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my error message: A")
}
