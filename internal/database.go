package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	Config StorageConfig
	db     *sql.DB
}

func NewDatabase(config StorageConfig) *Database {
	return &Database{config, nil}
}

func (conn *Database) MustConnect() {
	config := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conn.Config.User, conn.Config.Password, conn.Config.Database,
		conn.Config.Host, conn.Config.Port)
	db, err := sql.Open("postgres", config)
	if err != nil {
		panic(err)
	}

	conn.db = db
}

func (conn *Database) Disconnect() {
	conn.db.Close()
}
