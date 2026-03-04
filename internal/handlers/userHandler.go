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

type Mux interface {
	Handle(endpoint any, h tele.HandlerFunc, m ...tele.MiddlewareFunc)
	Use(middleware ...tele.MiddlewareFunc)
}

type UserHandler struct {
	userService *services.UserService
	linkService *services.LinkService
	conn        *pgx.Conn
}

func NewUserHandler(userService *services.UserService, linkService *services.LinkService, conn *pgx.Conn) *UserHandler {
	return &UserHandler{userService: userService, linkService: linkService, conn: conn}
}

func (h UserHandler) RegisterUser(ctx tele.Context) error {
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

func (h UserHandler) GetUser(ctx tele.Context) error {
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

func (h UserHandler) DeleteUser(c tele.Context) error {
	senderID := middleware.GetUserID(c)

	ctx := context.Background()

	tx, err := h.conn.Begin(ctx)

	if err != nil {
		slog.Error("UserHandler.DeleteUser creating tx", "err", err)

		return c.Send("some error occured")
	}

	defer tx.Rollback(ctx)

	userService := h.userService.NewWithTx(tx)
	linkService := h.linkService.NewWithTx(tx)

	err = linkService.DeleteUserLinks(ctx, senderID)

	if err != nil {
		slog.Error("UserHandler.DeleteUser deleting user links", "err", err, "telegram_id", senderID)

		return c.Send("some error occured")
	}

	err = userService.DeleteUser(ctx, senderID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Send(err.Error())
		}

		slog.Error("UserHandler.DeleteUser", "err", err)

		return c.Send("some error occured")
	}

	_ = tx.Commit(ctx)

	return c.Send("user deleted successfully")
}

func (h UserHandler) RegisterHandlers(mux Mux) {
	mux.Handle("/me", h.GetUser)
	mux.Handle("/stop", h.DeleteUser)
}
