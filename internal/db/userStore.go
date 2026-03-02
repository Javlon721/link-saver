package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/types"
	"github.com/jackc/pgx/v5"
)

type PostgreUserStore struct {
	conn  *pgx.Conn
	table string
}

func (p PostgreUserStore) GetUser(ctx context.Context, telegramID int64) (*types.User, error) {
	var user types.User

	query := fmt.Sprintf("select id, telegram_id from %s where telegram_id = $1", p.table)

	if err := p.conn.QueryRow(ctx, query, telegramID).Scan(&user.ID, &user.TelegramID); err != nil {

		slog.Error(err.Error())

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (p PostgreUserStore) Register(ctx context.Context, params *types.RegisterUser) (*types.User, error) {
	_, err := p.GetUser(ctx, params.TelegramID)

	if err == nil {
		return nil, errs.ErrUserAlreadyExists
	}

	if !errors.Is(err, errs.ErrUserNotFound) {
		return nil, err
	}

	query := fmt.Sprintf("insert into %s (telegram_id) values ($1) returning id", p.table)

	newUser := &types.User{
		TelegramID: params.TelegramID,
	}

	if err = p.conn.QueryRow(ctx, query, params.TelegramID).Scan(&newUser.ID); err != nil {
		slog.Error(err.Error())

		return nil, fmt.Errorf("some error occured while creating user")
	}

	return newUser, nil
}

func (p PostgreUserStore) DeleteUser(ctx context.Context, telegram_id int64) error {
	user, err := p.GetUser(ctx, telegram_id)

	if err != nil {
		return err
	}

	query := fmt.Sprintf("delete from %s where telegram_id = $1", p.table)

	_, err = p.conn.Exec(ctx, query, user.TelegramID)

	return err
}

func NewPostgresUserStore(conn *pgx.Conn, table string) *PostgreUserStore {
	cleanedTable := pgx.Identifier{table}.Sanitize()

	return &PostgreUserStore{
		conn:  conn,
		table: cleanedTable,
	}
}
