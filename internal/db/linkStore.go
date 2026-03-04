package db

import (
	"context"
	"fmt"

	"github.com/Javlon721/link-saver/internal/types"
	"github.com/jackc/pgx/v5"
)

type PostgreLinkStore struct {
	db    DB
	table string
}

func (store PostgreLinkStore) AddLink(ctx context.Context, userID int64, linkName, describtion string) (*types.Link, error) {
	query := fmt.Sprintf("insert into %s (user_id, link, describtion) values ($1, $2, $3) returning id", store.table)

	var link types.Link

	err := store.db.QueryRow(ctx, query, userID, linkName, describtion).Scan(&link.ID)

	if err != nil {
		return nil, err
	}

	link.Link = linkName
	link.Desctibtion = describtion
	link.UserID = userID

	return &link, err
}

func (store PostgreLinkStore) GetAll(ctx context.Context, user_id int64) []*types.Link {
	query := fmt.Sprintf("select id, user_id, link, describtion from %s where user_id = $1", store.table)

	var result []*types.Link
	rows, err := store.db.Query(ctx, query, user_id)

	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		link := &types.Link{}

		_ = rows.Scan(&link.ID, &link.UserID, &link.Link, &link.Desctibtion)

		result = append(result, link)
	}

	return result
}

func (store PostgreLinkStore) DeleteLink(ctx context.Context, userID, link_id int64) error {

	query := fmt.Sprintf("delete from %s where id = $1 and user_id = $2", store.table)

	_, err := store.db.Exec(ctx, query, link_id, userID)

	return err
}

func NewPostgreLinkStore(db DB, table string) *PostgreLinkStore {
	cleanedTable := pgx.Identifier{table}.Sanitize()

	return &PostgreLinkStore{
		db:    db,
		table: cleanedTable,
	}
}
