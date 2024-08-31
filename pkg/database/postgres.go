package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(host string, port int, user string, password string, database string, sslmode string) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, database, sslmode),
	)
	if err != nil {
		err = fmt.Errorf("sqlx.Connect(...): %w", err)
		return nil, err
	}
	return db, nil
}
