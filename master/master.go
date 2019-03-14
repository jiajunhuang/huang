package master

import (
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func Main(c *cli.Context) error {
	defer logger.Sync() // flushes buffer, if any

	return nil
}
