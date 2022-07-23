package telegram

import (
	"strconv"

	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/pkg/logging"
	"go.uber.org/fx"
	tele "gopkg.in/telebot.v3"
)

type InlineResultHandler struct {
	fx.In
	Logger  logging.Logger
	Usecase application.DownloadTorrentUsecase
}

func (h *InlineResultHandler) Handle(c tele.Context) error {
	h.Logger.Debug("chosen_inline_result", c.InlineResult().ResultID)

	torrentID, err := strconv.Atoi(c.InlineResult().ResultID)
	if err != nil {
		return err
	}

	notifier := NewTelegramNotifier(c)
	return h.Usecase.Execute(torrentID, notifier)
}
