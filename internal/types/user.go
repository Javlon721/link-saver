package types

import "context"

type User struct {
	ID         int64
	TelegramID int64
}

type RegisterUser struct {
	TelegramID int64
}

type UserStore interface {
	GetUser(context.Context, int64) (*User, error)
	Register(context.Context, *RegisterUser) (*User, error)
	DeleteUser(context.Context, int64) error
}
