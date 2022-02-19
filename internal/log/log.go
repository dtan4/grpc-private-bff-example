package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

func NewLogger() (*zap.Logger, error) {
	ec := zap.NewProductionEncoderConfig()

	// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
	ec.LevelKey = "severity"
	ec.EncodeLevel = encoderLevel

	// https://cloud.google.com/logging/docs/agent/logging/configuration#timestamp-processing
	ec.TimeKey = "time"
	ec.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	// https://cloud.google.com/logging/docs/agent/logging/configuration#special-fields
	ec.MessageKey = "message"

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = ec
	cfg.Sampling = nil
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	l, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return l, nil
}

func encoderLevel(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {
	pae.AppendString(logLevelSeverity[l])
}
