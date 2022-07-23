package main

import (
	"github.com/asyncee/homebot/internal/config"
	"github.com/asyncee/homebot/internal/telegram"
	"github.com/asyncee/homebot/internal/torrentsinfra"
	"github.com/asyncee/homebot/pkg/logging"
	"github.com/asyncee/homebot/pkg/rutracker"
	"github.com/asyncee/homebot/pkg/transmission"
	"go.uber.org/fx"
)

func main() {
	logging.Setup(logging.Debug)

	app := fx.New(
		telegram.Module,
		rutracker.Module,
		transmission.Module,
		config.Module,
		logging.Module,
		torrentsinfra.Module,
	)
	app.Run()
}
