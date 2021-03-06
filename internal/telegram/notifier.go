package telegram

import (
	"fmt"

	"github.com/asyncee/homebot/internal/torrents/application"
	tele "gopkg.in/telebot.v3"
)

type notifier struct {
	ctx tele.Context
}

func NewTelegramNotifier(ctx tele.Context) application.Notifier {
	return &notifier{ctx: ctx}
}

func (no *notifier) Notify(notification application.Notification) error {
	if notification.Link == nil {
		return no.ctx.Send(notification.Text)
	}
	markup := &tele.ReplyMarkup{}
	markup.Inline(
		markup.Row(
			markup.URL(notification.Link.Text, notification.Link.Url),
		),
	)
	return no.ctx.Send(
		notification.Text,
		markup,
	)
}

func (no *notifier) NotifyText(text string, args ...interface{}) error {
	return no.Notify(application.Notification{Text: fmt.Sprintf(text, args...)})
}

func (no *notifier) NotifyLink(text string, link *application.Link) error {
	return no.Notify(application.Notification{Text: text, Link: link})
}
