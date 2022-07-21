package main

import (
	"github.com/asyncee/homebot/internal/config"
	"github.com/asyncee/homebot/internal/telegram"
	"github.com/asyncee/homebot/internal/torrentsinfra"
	"github.com/asyncee/homebot/pkg/logging"
	"github.com/asyncee/homebot/pkg/rutracker"
	"go.uber.org/fx"
)

func main() {
	logging.Setup(logging.Debug)

	app := fx.New(
		telegram.Module,
		rutracker.Module,
		config.Module,
		logging.Module,
		torrentsinfra.Module,
		fx.Provide(
			func(cfg config.Config) (telegram.BotToken, telegram.TelegramAdmins) {
				return cfg.Telegram.Token, cfg.Telegram.Admins
			},
			func(cfg config.Config) (rutracker.Username, rutracker.Password) {
				return cfg.Rutracker.Login, cfg.Rutracker.Password
			},
		),
	)
	app.Run()
}
