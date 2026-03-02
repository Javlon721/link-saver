package templates

import (
	"bytes"
	"fmt"

	"github.com/Javlon721/link-saver/internal/types"
	tele "gopkg.in/telebot.v4"
)

func LinkMessageTemplate(link *types.Link) types.Message {
	return types.Message{
		Text:      linkTemplate(link),
		ParseMode: tele.ModeHTML,
	}
}

func linkTemplate(link *types.Link) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "<b>link</b>: %s\n", link.Link)

	if link.Desctibtion != "" {
		fmt.Fprintf(&buf, "<b>describtion</b>: %s", link.Desctibtion)
	}

	return buf.String()
}

func LinksTemplate(links []*types.Link) types.Message {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "<b>Links:</b>\n\n")

	for _, link := range links {
		fmt.Fprintf(&buf, "%s\n", linkTemplate(link))
	}

	return types.Message{
		Text:      buf.String(),
		ParseMode: tele.ModeHTML,
	}
}
