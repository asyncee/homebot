package telegram

import (
	"time"

	"go.uber.org/fx"
	tele "gopkg.in/telebot.v3"
)

type Handlers struct {
	fx.In
	OnSearchText OnSearchHandler
}

type Bot struct {
	bot *tele.Bot
}

func (b *Bot) Run() {
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

func (b *Bot) SetupHandlers(handlers Handlers) {
	// TODO: admin middleware

	b.bot.Handle("/start", onStartCommand)
	b.bot.Handle(tele.OnText, handlers.OnSearchText.Handle)
}

type BotToken string

type BotParams struct {
	fx.In
	Token    BotToken
	Handlers Handlers
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

	b := Bot{bot: telebot}
	b.SetupHandlers(p.Handlers)
	return &b, nil
}
