package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/zentrope/webl/internal"
)

type Database struct {
	Config internal.StorageConfig
	db     *sql.DB
}

func NewDatabase(config internal.StorageConfig) *Database {
	return &Database{config, nil}
}

func (conn *Database) Connect() {
	config := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conn.Config.User, conn.Config.Password, conn.Config.Database,
		conn.Config.Host, conn.Config.Port)
	db, err := sql.Open("postgres", config)
	if err != nil {
		log.Fatal(err)
	}

	conn.db = db
}

func (conn *Database) Disconnect() {
	conn.db.Close()
}

func (conn *Database) enums(typeName string) []string {

	rows, err := conn.db.Query("select enumlabel from vw_enums where typname=$1", typeName)
	if err != nil {
		log.Fatal(err)
	}

	results := make([]string, 0)

	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, s)
	}

	return results
}
