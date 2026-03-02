package db

import (
	"context"
	"fmt"

	"github.com/Javlon721/link-saver/internal/types"
	"github.com/jackc/pgx/v5"
)

type PostgreLinkStore struct {
	conn  *pgx.Conn
	table string
}

func (store PostgreLinkStore) Register(ctx context.Context, params *types.RegisterLink) (*types.Link, error) {
	query := fmt.Sprintf("insert into %s (user_id, link, describtion) values ($1, $2, $3) returning id", store.table)

	var link types.Link

	err := store.conn.QueryRow(ctx, query, params.UserID, params.Link, params.Desctibtion).Scan(&link.ID)

	if err != nil {
		return nil, err
	}

	link.Link = params.Link
	link.Desctibtion = params.Desctibtion
	link.UserID = params.UserID

	return &link, err
}

func (store PostgreLinkStore) GetAll(ctx context.Context) ([]*types.Link, error) {
	return nil, nil
}

func NewPostgreLinkStore(conn *pgx.Conn, table string) *PostgreLinkStore {
	cleanedTable := pgx.Identifier{table}.Sanitize()

	return &PostgreLinkStore{
		conn:  conn,
		table: cleanedTable,
	}
}
