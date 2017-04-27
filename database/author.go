package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Author struct {
	Handle string
	Email  string
	Type   string
	Status string
}

func (conn *Database) Authentic(handle, password string) bool {

	const q = "select handle from author where lower(handle)=lower($1) and password=$2"

	rows, err := conn.db.Query(q, handle, password)

	if err != nil {
		log.Fatal(err)
	}

	return rows.Next()
}

func (conn *Database) AuthorTypes() []string {
	return conn.enums("author_type")
}

func (conn *Database) Author(handle string) *Author {
	const query = "select handle, email, type, status from author where lower(handle)=lower($1)"
	rows, err := conn.db.Query(query, handle)

	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	rows.Next()
	return rowToAuthor(rows)
}

func (conn *Database) AuthorExists(handle string) bool {
	const query = "select handle from author where lower(handle)=lower($1)"
	rows, err := conn.db.Query(query, handle)

	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	return rows.Next()
}

func (conn *Database) Authors() []*Author {
	rows, err := conn.db.Query("select handle, email, type, status from author")
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	authors := make([]*Author, 0)

	for rows.Next() {
		authors = append(authors, rowToAuthor(rows))
	}

	return authors
}

func (conn *Database) CreateAuthor(handle, email, password string) error {
	_, err := conn.db.Exec(
		"insert into author (handle, email, password) values ($1, $2, $3)",
		handle, email, password,
	)

	return err
}

func (conn *Database) DeleteAuthor(handle string) {
	_, err := conn.db.Exec("delete from author where lower(handle)=lower($1)", handle)
	if err != nil {
		log.Fatal(err)
	}
}

func rowToAuthor(rows *sql.Rows) *Author {
	var a Author
	err := rows.Scan(&a.Handle, &a.Email, &a.Status, &a.Type)
	if err != nil {
		log.Fatal(err)
	}
	return &a
}
