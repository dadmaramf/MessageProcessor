package postgres

import (
	"database/sql"
	"messageprocessor/internal/config"
)

func NewPostgresDB(cfg config.ConfigInterface) (*sql.DB, error) {
	connStr := "user=" + cfg.GetPostgresConnect().User +
		" dbname=" + cfg.GetPostgresConnect().DBname +
		" password=" + cfg.GetPostgresConnect().Password +
		" sslmode=disable"

	con, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return con, nil
}
