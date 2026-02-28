package templates

import (
	"fmt"

	"github.com/Javlon721/link-saver/internal/types"
	tele "gopkg.in/telebot.v4"
)

func HelpDesk() types.Message {
	return types.Message{
		Text: fmt.Sprint(
			"<b>What can this bot do:</b>\n\n",
			"/help - prints help desk\n\n",
			"/links - print all links that you send\n\n",
		),
		ParseMode: tele.ModeHTML,
	}
}
