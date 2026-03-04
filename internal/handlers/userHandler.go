package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/middleware"
	"github.com/Javlon721/link-saver/internal/services"
	"github.com/Javlon721/link-saver/internal/types"
	"github.com/jackc/pgx/v5"
	tele "gopkg.in/telebot.v4"
)

type UserHandler struct {
	userService *services.UserService
	linkService *services.LinkService
}

func NewUserHandler(userService *services.UserService, linkService *services.LinkService) *UserHandler {
	return &UserHandler{userService: userService, linkService: linkService}
}

func (h *UserHandler) RegisterUser(ctx tele.Context) error {
	senderID := ctx.Sender().ID

	params := &types.RegisterUser{TelegramID: senderID}

	_, err := h.userService.RegisterUser(context.Background(), params)

	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return ctx.Send(err.Error())
		}

		slog.Error("userHandler.RegisterUser", "err", err)

		return nil
	}

	return nil
}

func (h *UserHandler) GetUser(ctx tele.Context) error {
	senderID := ctx.Sender().ID

	user, err := h.userService.GetUser(context.Background(), senderID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return ctx.Send("You need to register first")
		}

		slog.Error("userHandler.GetUser", "err", err)

		return nil
	}

	return ctx.Send(fmt.Sprintf("your id is: %d", user.TelegramID))
}

func (h *UserHandler) DeleteUser(c tele.Context) error {
	senderID := middleware.GetUserID(c)

	ctx, err := middleware.GetContext(c)
	if err != nil {
		slog.Error("UserHandler.DeleteUser getting ctx", "err", err)

		return c.Send("some error occured")
	}

	tx, err := middleware.GetTran(c)
	if err != nil {
		slog.Error("UserHandler.DeleteUser getting tx", "err", err)

		return c.Send("some error occured")
	}

	h = h.newWithTx(tx)

	err = h.linkService.DeleteUserLinks(ctx, senderID)

	if err != nil {
		slog.Error("UserHandler.DeleteUser deleting user links", "err", err, "telegram_id", senderID)

		return c.Send("some error occured")
	}

	err = h.userService.DeleteUser(ctx, senderID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Send(err.Error())
		}

		slog.Error("UserHandler.DeleteUser", "err", err)

		return c.Send("some error occured")
	}

	return c.Send("user deleted successfully")
}

func (h *UserHandler) newWithTx(tx pgx.Tx) *UserHandler {
	userService := h.userService.NewWithTx(tx)
	linkService := h.linkService.NewWithTx(tx)

	return &UserHandler{
		userService: userService,
		linkService: linkService,
	}
}
