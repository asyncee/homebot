package main

import (
	"github.com/asyncee/homebot/internal/config"
	"github.com/asyncee/homebot/internal/logging"
)

func main() {
	logging.Setup(logging.Debug)
	logger := logging.GetLogger()

	logger.Debug("msg", "initializing configuration")
	cfg := config.NewConfig()
	logger.Info("config", cfg)
}
