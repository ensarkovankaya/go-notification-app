package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// init initializes the logger
func init() {
	var err error
	cfg := zap.NewProductionConfig()
	loglevel, _ := zapcore.ParseLevel(Cnf.LogLevel)
	cfg.Level = zap.NewAtomicLevelAt(loglevel)
	if logger, err = cfg.Build(); err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	zap.ReplaceGlobals(logger)
}
