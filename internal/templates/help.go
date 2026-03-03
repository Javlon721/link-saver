package templates

import (
	"bytes"
	"fmt"

	"github.com/Javlon721/link-saver/internal/types"
	tele "gopkg.in/telebot.v4"
)

func HelpDesk() types.Message {
	var buf bytes.Buffer

	// essentials
	fmt.Fprint(&buf, fmt.Sprint(
		"<b>What can this bot do:</b>\n",
		"/help - prints help desk\n",
		"\n",
	))

	// all user story
	fmt.Fprint(&buf, fmt.Sprint(
		"<b>Your account info:</b>\n",
		"/register - registers new user\n",
		"/me - send current user info\n",
		"/stop - deletes all your data (<b>DANGER</b>)\n",
		"\n",
	))

	// all links story
	fmt.Fprint(&buf, fmt.Sprint(
		"<b>Your links info:</b>\n",
		"/link - register new link\n",
		"/links - print all links that you have\n",
		"/linksBtns - print all links that you have + alter buttons\n",
		"\n",
	))

	return types.Message{
		Text:      buf.String(),
		ParseMode: tele.ModeHTML,
	}
}
