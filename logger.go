package main

import (
	"fmt"
	"go.uber.org/zap"
)

var logger *zap.Logger

// init initializes the logger
func init() {
	var err error
	if logger, err = zap.NewProduction(); err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	zap.ReplaceGlobals(logger)
}
