package handlers

import (
	"github.com/Javlon721/link-saver/internal/templates"
	tele "gopkg.in/telebot.v4"
)

type MainHandler struct{}

func (h MainHandler) HelpDeskHandler(ctx tele.Context) error {
	message := templates.HelpDesk()

	return ctx.Send(message.Text, message.ParseMode)
}

func (h MainHandler) RegisterHandlers(bot *tele.Bot) {
	bot.Handle("/start", h.HelpDeskHandler)
	bot.Handle("/help", h.HelpDeskHandler)
}

func NewMainHandler() *MainHandler {
	return &MainHandler{}
}
