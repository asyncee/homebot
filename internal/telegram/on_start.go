package telegram

import (
	"github.com/asyncee/homebot/pkg/logging"
	"go.uber.org/fx"
	tele "gopkg.in/telebot.v3"
)

type StartCommandHandler struct {
	fx.In
	Logger logging.Logger
}

func (h *StartCommandHandler) Handle(c tele.Context) error {
	h.Logger.Debug("command", "/start")
	return c.Send("Привет! Что мне найти для тебя?")
}
