package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const logPath = "./log/out.log"

// NewLoggerOrFail returns a new logger or panics.
func NewLoggerOrFail() (*zap.Logger, func() error) {
	if _, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0600); err != nil {
		panic(err)
	}

	// log retention policy
	lj := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1000,
		MaxBackups: 100,
		MaxAge:     30,
		Compress:   true,
	}

	ws := zapcore.AddSync(lj)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, ws, zapcore.DebugLevel),
		zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, logger.Sync
}
