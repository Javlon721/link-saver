package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/services"
	"github.com/Javlon721/link-saver/internal/templates"
	"github.com/Javlon721/link-saver/internal/types"
	tele "gopkg.in/telebot.v4"
)

type LinkHandler struct {
	linkService *services.LinkService
}

type Sendable interface {
	Send(to tele.Recipient, what any, opts ...any) (*tele.Message, error)
}

type CallbackRegistrer interface {
	HandlerCallback(string, tele.HandlerFunc)
}

func NewLinkHandler(linkService *services.LinkService) *LinkHandler {
	return &LinkHandler{linkService: linkService}
}

func (h LinkHandler) RegisterLink(c tele.Context) error {
	senderID := c.Sender().ID

	ctx := context.Background()

	payload := strings.TrimSpace(c.Message().Payload)

	if payload == "" {
		return c.Send("you need to provide link")
	}

	args := strings.SplitN(payload, " ", 2)

	params := &types.RegisterLink{
		UserID: senderID,
		Link:   args[0],
	}

	if len(args) == 2 {
		params.Desctibtion = args[1]
	}

	link, err := h.linkService.RegisterLink(ctx, params)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Send("You need to register first")
		}

		slog.Error("LinkHandler.RegisterLink", "err", err)

		return nil
	}

	return c.Send(fmt.Sprintf("link registered: %d", link.ID))
}

func (h LinkHandler) GetAll(c tele.Context) error {
	senderID := c.Sender().ID

	ctx := context.Background()

	links, err := h.linkService.GetAll(ctx, senderID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Send("You need to register first")
		}

		if errors.Is(err, errs.ErrLinksNotFound) {
			return c.Send("You out of links")
		}

		slog.Error("LinkHandler.GetAll", "err", err)

		return nil
	}

	message := templates.LinksTemplate(links)

	return c.Send(message.Text, message.ParseMode)
}

func (h LinkHandler) GetAllWithBtns(c tele.Context) error {
	sender := c.Sender()

	ctx := context.Background()

	links, err := h.linkService.GetAll(ctx, sender.ID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Send("You need to register first")
		}

		if errors.Is(err, errs.ErrLinksNotFound) {
			return c.Send("You out of links")
		}

		slog.Error("LinkHandler.GetAll", "err", err)

		return nil
	}

	for _, link := range links {
		go func(sender tele.Recipient, bot Sendable, link *types.Link) {

			menu := &tele.ReplyMarkup{}

			message, btn := templates.LinkTemplateWithBtns(link, menu)

			menu.Inline(menu.Row(*btn))

			bot.Send(sender, message.Text, message.ParseMode, menu)

		}(sender, c.Bot(), link)
	}

	return nil
}

func (h LinkHandler) RegisterHandlers(bot *tele.Bot) {
	bot.Handle("/link", h.RegisterLink)
	bot.Handle("/links", h.GetAll)
	bot.Handle("/linksBtns", h.GetAllWithBtns)
}

func (h LinkHandler) DeleteLink(c tele.Context) error {
	payload := strings.SplitN(c.Callback().Data, "|", 2)

	if len(payload) < 2 {
		slog.Error("LinkHandler.DeleteLink", "err", "need to privide link id for deletion", "payload", payload)
		return nil
	}

	ctx := context.Background()

	linkID, err := strconv.ParseInt(payload[1], 10, 64)

	if err != nil {
		slog.Error("LinkHandler.DeleteLink", "err", err)
		return nil
	}

	err = h.linkService.DeleteLink(ctx, c.Sender().ID, linkID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Send("You need to register first")
		}

		slog.Error("LinkHandler.DeleteLink", "err", err)

		return nil
	}

	return c.Delete()
}

func (h LinkHandler) GetCallbacks() map[string]tele.HandlerFunc {
	callbacks := map[string]tele.HandlerFunc{}

	callbacks["delete"] = h.DeleteLink

	return callbacks
}
