package zlog

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"strconv"
)

var logger *zap.Logger

func InitDefaultZapLog(fileName string) {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	  	"messageKey": "message",
	  	"levelKey": "level",
		"timeKey": "time",
	  	"levelEncoder": "lowercase",
		"timeEncoder": "rfc3339"
	  }}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.OutputPaths = append(cfg.OutputPaths, fileName)
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = log
	logger.Info("logger construction succeeded")
}
func InitAdvanceZapLog() {
	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	topicDebugging := zapcore.AddSync(ioutil.Discard)
	topicErrors := zapcore.AddSync(ioutil.Discard)
	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	logger.Info("Constructed a advance logger")
}
func Debug(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	}
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}
func Uint16(key string, val uint16) zap.Field {
	return zap.Uint16(key, val)
}
func ByteString(key string, val []byte) zap.Field {
	return zap.ByteString(key, val)
}
func Uint32(key string, val uint32) zap.Field {
	return zap.Uint32(key, val)
}
func Int32(key string, val int32) zap.Field {
	return zap.Int32(key, val)
}
func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}
func Uint32Arr(key string, val []uint32) zap.Field {
	var str bytes.Buffer
	for _, value := range val {
		str.WriteString(strconv.FormatUint(uint64(value), 10))
	}
	return zap.String(key, str.String())
}
