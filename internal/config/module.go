package config

import (
	"time"

	"github.com/asyncee/homebot/internal/telegram"
	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/internal/torrentsinfra"
	"github.com/asyncee/homebot/pkg/rutracker"
	"github.com/asyncee/homebot/pkg/transmission"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewConfig,
		func(cfg Config) (telegram.BotToken, telegram.TelegramAdmins) {
			return cfg.Telegram.Token, cfg.Telegram.Admins
		},
		func(cfg Config) (rutracker.Username, rutracker.Password) {
			return cfg.Rutracker.Login, cfg.Rutracker.Password
		},
		func(cfg Config) torrentsinfra.DownloadTorrentsDir {
			return cfg.App.DownloadTorrentsTo
		},
		func(cfg Config) application.ViewTorrentInUILink {
			return application.ViewTorrentInUILink(cfg.Transmission.WebUiLink)
		},
		func(cfg Config) (torrentsinfra.PollStatusDuration, torrentsinfra.PollStatusTimeout) {
			duration := torrentsinfra.PollStatusDuration(time.Second * time.Duration(cfg.App.PollTorrentStatusDurationSeconds))
			timeout := torrentsinfra.PollStatusTimeout(time.Second * time.Duration(cfg.App.PollTorrentStatusTimeoutSeconds))
			return duration, timeout
		},
		func(cfg Config) (transmission.Host, transmission.User, transmission.Password) {
			host := transmission.Host(cfg.Transmission.RPCHost)
			user := transmission.User(cfg.Transmission.RPCUser)
			password := transmission.Password(cfg.Transmission.RPCPassword)
			return host, user, password
		},
	),
)
