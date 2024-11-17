package logger

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	envLogLevel      = "LOG_LEVEL"
	envLogOutput     = "LOG_OUTPUT"
	logListMaxLength = 700
	logListTrimBy    = 100
)

var (
	log     logger
	sink    *MemorySink
	loglist []LogEntry
)

type loggerInterface interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Write([]byte) (int, error)
}

type logger struct {
	log *zap.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

type MemorySink struct {
	*bytes.Buffer
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

type LogEntry struct {
	LogTime    string
	LogLevel   string
	LogMessage string
}

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

func init() {
	initLogger(false, "")
}

func Init(logFileName string) {
	initLogger(false, logFileName)
}

func initLogger(test bool, logFileName string) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "msg"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.StacktraceKey = ""

	if logFileName != "" {
		zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
			return lumberjackSink{
				Logger: &lumberjack.Logger{
					Filename:   u.Opaque,
					MaxSize:    100, // MB
					MaxAge:     7,   // days
					MaxBackups: 7,
					LocalTime:  false,
					Compress:   true,
				},
			}, nil
		})
	}

	logConfig := zap.NewProductionConfig()
	if test {
		zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
			return sink, nil
		})
		sink = &MemorySink{new(bytes.Buffer)}
		logConfig.OutputPaths = []string{"memory://"}
	} else {
		logConfig.OutputPaths = []string{getOutput()}
	}
	if logFileName != "" {
		logConfig.OutputPaths = append(logConfig.OutputPaths, fmt.Sprintf("lumberjack:%v", logFileName))
	}

	logConfig.Level = zap.NewAtomicLevelAt(getLevel())
	logConfig.Encoding = "json"
	logConfig.EncoderConfig = encoderConfig

	var err error
	if log.log, err = logConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1)); err != nil {
		panic(err)
	}
}

func getLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	output := strings.TrimSpace(os.Getenv(envLogOutput))
	if output == "" {
		return "stdout"
	}
	return output
}

func GetLogger() loggerInterface {
	return log
}

func (l logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v...))
	}
}

func (l logger) Print(v ...interface{}) {
	Info(fmt.Sprintf("%v", v))
}

func (l logger) Write(data []byte) (n int, err error) {
	logMessage := string(data)
	if strings.Contains(strings.ToLower(logMessage), "error") {
		Error(logMessage, nil)
	} else if strings.Contains(strings.ToLower(logMessage), "warn") {
		Warn(logMessage)
	} else if strings.Contains(strings.ToLower(logMessage), "debug") {
		Debug(logMessage)
	} else {
		Info(logMessage)
	}
	return len(data), nil
}

func fieldsToZapField(tags []Field) []zapcore.Field {
	zapTags := make([]zap.Field, 0)
	for _, tag := range tags {
		zapTag := zap.Any(tag.Key, tag.Value)
		zapTags = append(zapTags, zapTag)
	}
	return zapTags
}

func trimList() {
	if len(loglist) > logListMaxLength {
		loglist = loglist[logListTrimBy:]
	}
}

func addToLogList(logLevel string, msg string) {
	var entry LogEntry
	entry.LogTime = time.Now().Format(time.RFC3339)
	entry.LogLevel = logLevel
	entry.LogMessage = msg
	loglist = append(loglist, entry)
	trimList()
}

func GetLogList() []LogEntry {
	return loglist
}

func ClearLogList() {
	loglist = loglist[:0]
}

func Debugf(msg string, a ...any) {
	m := fmt.Sprintf(msg, a...)
	addToLogList("Debug", m)
	log.log.Debug(m)
	log.log.Sync()
}

func Debug(msg string, tags ...Field) {
	zapTags := fieldsToZapField(tags)
	log.log.Debug(msg, zapTags...)
	log.log.Sync()
}

func Infof(msg string, a ...any) {
	m := fmt.Sprintf(msg, a...)
	addToLogList("Info", m)
	log.log.Info(m)
	log.log.Sync()
}

func Info(msg string, tags ...Field) {
	addToLogList("Info", msg)
	zapTags := fieldsToZapField(tags)
	log.log.Info(msg, zapTags...)
	log.log.Sync()
}

func Warnf(msg string, a ...any) {
	m := fmt.Sprintf(msg, a...)
	addToLogList("Warn", m)
	log.log.Warn(m)
	log.log.Sync()
}

func Warn(msg string, tags ...Field) {
	addToLogList("Warn", msg)
	zapTags := fieldsToZapField(tags)
	log.log.Warn(msg, zapTags...)
	log.log.Sync()
}

func Errorf(msg string, a ...any) {
	m := fmt.Sprintf(msg, a...)
	addToLogList("Error", m)
	log.log.Error(m)
	log.log.Sync()
}

func Error(msg string, err error, tags ...Field) {
	var m string
	if err != nil {
		m = msg + ": " + err.Error()
	} else {
		m = msg
	}
	addToLogList("Error", m)
	zapTags := fieldsToZapField(tags)
	zapTags = append(zapTags, zap.NamedError("error", err))
	log.log.Error(msg, zapTags...)
	log.log.Sync()
}
