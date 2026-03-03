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
	linkStore   types.LinkStore
	userStore   types.UserStore
	linkService *services.LinkService
}

type Sendable interface {
	Send(to tele.Recipient, what any, opts ...any) (*tele.Message, error)
}

type CallbackRegistrer interface {
	HandlerCallback(string, tele.HandlerFunc)
}

func NewLinkHandler(linkStore types.LinkStore, userStore types.UserStore, linkService *services.LinkService) *LinkHandler {
	return &LinkHandler{linkStore: linkStore, userStore: userStore, linkService: linkService}
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

func (h LinkHandler) GetAll(ctx tele.Context) error {
	senderID := ctx.Sender().ID

	contextBG := context.Background()

	user, err := h.userStore.GetUser(contextBG, senderID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return ctx.Send("You need to register first")
		}

		slog.Error("LinkHandler.GetAll", "err", err)

		return nil
	}

	links := h.linkStore.GetAll(contextBG, user.ID)

	if len(links) == 0 {
		return ctx.Send("you does not have any links")
	}

	message := templates.LinksTemplate(links)

	return ctx.Send(message.Text, message.ParseMode)
}

func (h LinkHandler) GetAllWithBtns(ctx tele.Context) error {
	sender := ctx.Sender()

	contextBG := context.Background()

	user, err := h.userStore.GetUser(contextBG, sender.ID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return ctx.Send("You need to register first")
		}

		slog.Error("LinkHandler.GetAll", "err", err)

		return nil
	}

	links := h.linkStore.GetAll(contextBG, user.ID)

	if len(links) == 0 {
		return ctx.Send("you does not have any links")
	}

	for _, link := range links {
		go func(sender tele.Recipient, bot Sendable, link *types.Link) {

			menu := &tele.ReplyMarkup{}

			message, btn := templates.LinkTemplateWithBtns(link, menu)

			menu.Inline(menu.Row(*btn))

			bot.Send(sender, message.Text, message.ParseMode, menu)

		}(sender, ctx.Bot(), link)
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
		slog.Error("LinkHandler.DeleteLink payload", "err", "need to privide link id for deletion", "payload", payload)
		return nil
	}

	ctx := context.Background()

	linkID, err := strconv.ParseInt(payload[1], 10, 64)

	if err != nil {
		slog.Error("LinkHandler.DeleteLink linkID", "err", err)
		return nil
	}

	err = h.linkStore.DeleteLink(ctx, linkID)

	if err != nil {
		slog.Error("LinkHandler.DeleteLink linkID", "err", err)
		return nil
	}

	return c.Delete()
}

func (h LinkHandler) GetCallbacks() map[string]tele.HandlerFunc {
	callbacks := map[string]tele.HandlerFunc{}

	callbacks["delete"] = h.DeleteLink

	return callbacks
}
