package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/types"
	tele "gopkg.in/telebot.v4"
)

type LinkHandler struct {
	linkStore types.LinkStore
	userStore types.UserStore
}

func NewLinkHandler(linkStore types.LinkStore, userStore types.UserStore) *LinkHandler {
	return &LinkHandler{linkStore: linkStore, userStore: userStore}
}

func (h LinkHandler) RegisterLink(ctx tele.Context) error {
	senderID := ctx.Sender().ID

	contextBG := context.Background()

	user, err := h.userStore.GetUser(contextBG, senderID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return ctx.Send("You need to register first")
		}

		slog.Error("LinkHandler.RegisterLink", "err", err)

		return nil
	}

	payload := strings.TrimSpace(ctx.Message().Payload)

	if payload == "" {
		return ctx.Send("you need to provide link")
	}

	args := strings.SplitN(payload, " ", 2)

	params := &types.RegisterLink{
		UserID: user.ID,
	}

	params.Link = args[0]

	if len(args) == 2 {
		params.Desctibtion = args[1]
	}

	link, err := h.linkStore.Register(contextBG, params)

	if err != nil {
		slog.Error("LinkHandler.RegisterLink", "err", err)

		return nil
	}

	return ctx.Send(fmt.Sprintf("link registered: %d", link.ID))
}

func (h LinkHandler) RegisterHandlers(bot *tele.Bot) {
	bot.Handle("/link", h.RegisterLink)
}
