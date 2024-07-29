package postgres

import (
	"database/sql"
	"fmt"
	"messageprocessor/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.ConfigInterface) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.GetPostgresConnect().User,
		cfg.GetPostgresConnect().Password,
		cfg.GetPostgresConnect().Host,
		cfg.GetPostgresConnect().Port,
		cfg.GetPostgresConnect().DBname,
	)
	con, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return con, nil
}
