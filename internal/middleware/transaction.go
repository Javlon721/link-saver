package middleware

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	tele "gopkg.in/telebot.v4"
)

var (
	txKey  = "tx"
	ctxKey = "ctx"
)

type Begginer interface {
	Begin(context.Context) (pgx.Tx, error)
}

func BeginCommitRollback(begginer Begginer) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := context.Background()
			tx, err := begginer.Begin(ctx)

			if err != nil {
				slog.Error("BEGIN TRANSACTION", "err", err)
				return nil
			}

			defer tx.Rollback(ctx)

			c.Set(txKey, tx)
			c.Set(ctxKey, ctx)

			resp := next(c)

			if resp != nil {
				return resp
			}

			if err = tx.Commit(ctx); err != nil {
				slog.Error("COMMIT TRANSACTION", "err", err)
			}

			return resp
		}
	}
}

func GetTran(c tele.Context) (pgx.Tx, error) {
	tx, ok := c.Get(txKey).(pgx.Tx)

	if !ok {
		return nil, fmt.Errorf("transaction not found in tele.Context")
	}

	return tx, nil
}

func GetContext(c tele.Context) (context.Context, error) {
	ctx, ok := c.Get(ctxKey).(context.Context)

	if !ok {
		return nil, fmt.Errorf("context not found in tele.Context")
	}

	return ctx, nil
}
