package handlers

import (
	"fmt"

	"github.com/Javlon721/link-saver/internal/templates"
	tele "gopkg.in/telebot.v4"
)

type MainHandler struct{}

func (h MainHandler) HelpDeskHandler(ctx tele.Context) error {
	message := templates.HelpDesk()

	fmt.Println("Message:", message.Text)

	return ctx.Send(message.Text, message.ParseMode)
}

func (h MainHandler) RegisterHandlers(bot *tele.Bot) {
	bot.Handle("/start", h.HelpDeskHandler)
	bot.Handle("/help", h.HelpDeskHandler)
}

func NewMainHandler() *MainHandler {
	return &MainHandler{}
}
