package telegram

import (
	"context"

	"github.com/asyncee/homebot/pkg/logging"
	"go.uber.org/fx"
)

func provideTelegramBot(
	lc fx.Lifecycle,
	logger logging.Logger,
	p BotParams,
) (*Bot, error) {
	bot, err := NewBot(p)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Debug("msg", "starting bot...")
			go bot.Run()
			logger.Info("msg", "bot started")
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Debug("msg", "stopping bot...")
			bot.Stop()
			logger.Info("msg", "bot stopped")
			return nil
		},
	})

	return bot, nil
}

var Module = fx.Options(
	fx.Provide(provideTelegramBot),
	fx.Invoke(func(b *Bot) {}),
)
