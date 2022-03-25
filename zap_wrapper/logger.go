package zap_wrapper

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(file string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.OutputPaths = []string{"stderr", file}
	config.ErrorOutputPaths = []string{"stderr", file}
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logger, _ := config.Build()
	return logger
}

var DefaultLogger = NewLogger("ddns.log")
