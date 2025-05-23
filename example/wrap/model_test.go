package wrap_test

import (
	"fmt"
	"testing"

	"github.com/AugustineAurelius/eos/example/wrap"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestTest1(t *testing.T) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	wrappedTest := wrap.NewTestMiddleware(&wrap.Test{}, wrap.WithTestLogging(logger), wrap.WithTestCircuitBreaker(wrap.NewCircuitBreakerConfig()))

	for range 6 {
		_, err = wrappedTest.Test1(1, &wrap.Test222{Name: "123"})
		fmt.Println(err)
	}
	require.Equal(t, err, wrap.ErrOpenCircuitBreaker)
	// require.NoError(t, err)
	// require.Equal(t, 13, res)

}
