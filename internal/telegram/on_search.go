package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/pkg/logging"
	"go.uber.org/fx"
	tele "gopkg.in/telebot.v3"
)

type TextHandler struct {
	fx.In
	Logger logging.Logger
	Repo   application.TorrentRepository
}

func (h *TextHandler) Handle(c tele.Context) error {
	query := strings.TrimSpace(strings.TrimPrefix(c.Text(), "@"+c.Bot().Me.Username))
	h.Logger.Debug("text", query)

	// TODO: extract TorrentsCountByNameQuery

	torrents, err := h.Repo.FindByName(context.TODO(), query)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка: %s", err.Error()))
	}

	if len(torrents) == 0 {
		return c.Send("По этому запросу я ничего не нашёл :(")
	}

	markup := &tele.ReplyMarkup{}
	markup.Inline(
		markup.Row(
			markup.QueryChat("Выбрать торрент", query),
		),
	)

	return c.Send(
		fmt.Sprintf("По запросу '%s' найдено торрентов: %d", query, len(torrents)),
		markup,
	)
}
