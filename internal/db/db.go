package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Javlon721/link-saver/internal/config"
	"github.com/jackc/pgx/v5"
)

func NewPostgreConn(config *config.Config) (*pgx.Conn, error) {
	uri := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.POSTGRES_USER, config.POSTGRES_PASSWORD),
		Host:   fmt.Sprintf("%s:%d", config.POSTGRES_HOST, config.POSTGRES_PORT),
		Path:   config.POSTGRES_DB,
	}

	conn, err := pgx.Connect(context.Background(), uri.String())

	return conn, err
}
