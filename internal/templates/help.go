package templates

import (
	"fmt"

	"github.com/Javlon721/link-saver/internal/types"
	tele "gopkg.in/telebot.v4"
)

func HelpDesk() types.Message {
	return types.Message{
		Text: fmt.Sprint(
			"<b>What can this bot do:</b>\n",
			"/register - registers new user\n",
			"/me - send current user info\n",
			"/help - prints help desk\n",
			"/links - print all links that you send\n",
		),
		ParseMode: tele.ModeHTML,
	}
}
