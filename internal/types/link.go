package types

import "context"

type LinkStore interface {
	Register(context.Context, *RegisterLink) (*Link, error)
	GetAll(context.Context, int64) []*Link
}

type RegisterLink struct {
	UserID      int64
	Link        string
	Desctibtion string
}

type Link struct {
	ID          int64
	UserID      int64
	Link        string
	Desctibtion string
}
