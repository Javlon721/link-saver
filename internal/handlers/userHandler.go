package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/services"
	"github.com/Javlon721/link-saver/internal/types"
	tele "gopkg.in/telebot.v4"
)

type UserHandler struct {
	userStore   types.UserStore
	userService *services.UserService
}

func NewUserHandler(userStore types.UserStore, userService *services.UserService) *UserHandler {
	return &UserHandler{userStore: userStore, userService: userService}
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

func (h UserHandler) DeleteUser(ctx tele.Context) error {
	senderID := ctx.Sender().ID

	err := h.userService.DeleteUser(context.Background(), senderID)

	if err != nil {

		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return ctx.Send(err.Error())
		}

		slog.Error("UserHandler.DeleteUser", "err", err)

		return ctx.Send("some error occured")
	}

	return ctx.Send("user deleted successfully")
}

func (h UserHandler) RegisterHandlers(bot *tele.Bot) {
	bot.Handle("/register", h.RegisterUser)
	bot.Handle("/me", h.GetUser)
	bot.Handle("/stop", h.DeleteUser)
}
