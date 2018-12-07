package worker

import (
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func Main() {
	defer logger.Sync() // flushes buffer, if any
}
