package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var m Moderer

type Logger struct {
	*zap.Logger
}

type Moderer interface {
	IsProduction() bool
}

func New[M Moderer](mode M, opts ...Opt) *Logger {
	m = mode
	var core zapcore.Core
	if mode.IsProduction() {
		core = productionMode()
	} else {
		core = developMode()
	}

	logger := zap.New(core)
	return &Logger{logger}
}

func productionMode() zapcore.Core {
	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	})

	return zapcore.NewCore(jsonEncoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)
}

func developMode() zapcore.Core {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentConfig().EncoderConfig)
	return zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
}

// Options
type Opt func() zapcore.Core

func WithFile(ws zapcore.WriteSyncer) zapcore.Core {
	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	})

	return zapcore.NewCore(jsonEncoder, zapcore.Lock(ws), zapcore.InfoLevel)
}
