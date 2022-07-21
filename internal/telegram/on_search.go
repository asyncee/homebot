package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/asyncee/homebot/internal/torrents/application"
	"go.uber.org/fx"
	tele "gopkg.in/telebot.v3"
)

type OnSearchHandler struct {
	fx.In
	Repo application.TorrentRepository
}

func (h *OnSearchHandler) Handle(c tele.Context) error {
	query := strings.TrimSpace(strings.TrimPrefix(c.Text(), "@"+c.Bot().Me.Username))

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
