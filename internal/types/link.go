package types

import "context"

type LinkStore interface {
	AddLink(ctx context.Context, userID int64, linkName, describtion string) (*Link, error)
	GetAll(context.Context, int64) []*Link
	DeleteLink(context.Context, int64, int64) error
	DeleteUserLinks(context.Context, int64) error
	NewWithTx(db DB) LinkStore
}

type LinkInfo struct {
	Link        string
	Desctibtion string
}

type Link struct {
	ID          int64
	UserID      int64
	Link        string
	Desctibtion string
}
