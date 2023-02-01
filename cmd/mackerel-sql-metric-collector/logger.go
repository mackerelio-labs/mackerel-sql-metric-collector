package main

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildLogger(logFormat, logLevel string) (logr.Logger, error) {
	var logger logr.Logger

	conf := zap.NewProductionConfig()
	conf.Encoding = logFormat
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var l zapcore.Level
	if err := l.UnmarshalText([]byte(logLevel)); err != nil {
		return logger, err
	}
	conf.Level = zap.NewAtomicLevelAt(l)

	zapLog, err := conf.Build()
	if err != nil {
		return logger, err
	}

	return zapr.NewLogger(zapLog), nil
}
