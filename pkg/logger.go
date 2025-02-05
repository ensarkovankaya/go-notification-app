package pkg

import (
	"fmt"
	"go.uber.org/zap"
)

var Logger *zap.Logger

// init initializes the logger
func init() {
	var err error
	if Logger, err = zap.NewProduction(); err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
}
