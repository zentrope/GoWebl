package webl

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	User     string
	Password string
	Database string
	db       *sql.DB
}

func NewDatabase(user, pass, dbname string) *Database {
	return &Database{user, pass, dbname, nil}
}

func (conn *Database) Connect() {
	config := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		conn.User, conn.Password, conn.Database)
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
