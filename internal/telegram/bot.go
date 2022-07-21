package telegram

import (
	"time"

	"go.uber.org/fx"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Handlers struct {
	fx.In
	OnSearchText         OnSearchHandler
	OnInlineQueryHandler OnInlineQueryHandler
}

type TelegramUserId int64
type TelegramAdmins []TelegramUserId

func (t TelegramAdmins) toInt64Slice() []int64 {
	result := make([]int64, len(t))
	for i := range t {
		result = append(result, int64(t[i]))
	}
	return result
}

type Bot struct {
	bot    *tele.Bot
	admins TelegramAdmins
}

func (b *Bot) Run() {
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

func (b *Bot) setupMiddleware() {
	b.bot.Use(middleware.IgnoreVia())
	b.bot.Use(middleware.Whitelist(b.admins.toInt64Slice()...))
}

func (b *Bot) setupHandlers(handlers Handlers) {
	b.bot.Handle("/start", onStartCommand)
	b.bot.Handle(tele.OnQuery, handlers.OnInlineQueryHandler.Handle)
	b.bot.Handle(tele.OnText, handlers.OnSearchText.Handle)
}

type BotToken string

type BotParams struct {
	fx.In
	Token    BotToken
	Handlers Handlers
	Admins   TelegramAdmins
}

func NewBot(p BotParams) (*Bot, error) {
	pref := tele.Settings{
		Token:  string(p.Token),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	telebot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b := Bot{bot: telebot, admins: p.Admins}
	b.setupMiddleware()
	b.setupHandlers(p.Handlers)
	return &b, nil
}
