package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Telegram struct {
		Token  string  `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
		Admins []int64 `env:"TELEGRAM_BOT_ADMINS" env-required:"true"`
	}
	Rutracker struct {
		Login    string `env:"RUTRACKER_LOGIN" env-required:"true"`
		Password string `env:"RUTRACKER_PASSWORD" env-required:"true"`
	}
	Transmission struct {
		RPCHost     string `env:"TRANSMISSION_RPC_HOST" env-default:"localhost"`
		RPCUser     string `env:"TRANSMISSION_RPC_USER" env-default:"admin"`
		RPCPassword string `env:"TRANSMISSION_RPC_PASSWORD" env-required:"true"`
		WebUiLink   string `env:"TRANSMISSION_WEB_UI_LINK" env-default:"http://localhost:9091/transmission/web"`
	}
	App struct {
		DownloadTorrentsTo               string `env:"DOWNLOAD_DIR" env-required:"true"`
		PollTorrentStatusTimeoutSeconds  int    `env:"POLL_TORRENT_STATUS_TIMEOUT_SECONDS" env-default:"3600"`
		PollTorrentStatusDurationSeconds int    `env:"POLL_TORRENT_STATUS_DURATION_SECONDS" env-default:"5"`
	}
}

func NewConfig() Config {
	cfg := Config{}
	cleanenv.ReadConfig(".env", &cfg)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		help, _ := cleanenv.GetDescription(&cfg, nil)
		log.Fatalf("failed to initialize configuration: %v\n%s", err, help)
	}
	return cfg
}
