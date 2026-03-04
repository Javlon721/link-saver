package middleware

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/services"
	tele "gopkg.in/telebot.v4"
)

func AuthorizeUser(userService *services.UserService) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := context.Background()

			senderID := c.Sender().ID

			user, err := userService.GetUser(ctx, senderID)

			if err != nil {
				if errors.Is(err, errs.ErrUserNotFound) {
					return c.Send("You need to register first")
				}

				slog.Error("AuthorizeUser", "err", err)

				return nil
			}

			c.Set("UserID", user.ID)

			return next(c)
		}
	}
}

func GetUserID(c tele.Context) int64 {
	userID, ok := c.Get("UserID").(int64)

	if !ok {
		panic("UserID inside of context should be int64")
	}

	return userID
}
