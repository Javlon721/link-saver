package db

import (
	"context"
	"fmt"

	"github.com/Javlon721/link-saver/internal/config"
	"github.com/jackc/pgx/v5"
)

func NewPostgreConn(config *config.Config) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.POSTGRES_USER, config.POSTGRES_PASSWORD, config.POSTGRES_HOST, config.POSTGRES_PORT, config.POSTGRES_DB)

	conn, err := pgx.Connect(context.Background(), connString)

	return conn, err
}
