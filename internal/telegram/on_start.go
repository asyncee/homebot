package telegram

import tele "gopkg.in/telebot.v3"

func onStartCommand(c tele.Context) error {
	return c.Send("Привет! Что мне найти для тебя?")
}
