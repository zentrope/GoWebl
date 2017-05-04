// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Author struct {
	Handle string
	Email  string
	Type   string
	Status string
}

func (conn *Database) Authentic(handle, password string) (bool, error) {

	const q = "select handle from author where lower(handle)=lower($1) and password=$2"

	rows, err := conn.db.Query(q, handle, password)

	defer rows.Close()

	if err != nil {
		return false, err
	}

	return rows.Next(), err
}

func (conn *Database) Author(handle string) (*Author, error) {
	const query = "select handle, email, type, status from author where lower(handle)=lower($1)"
	rows, err := conn.db.Query(query, handle)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	rows.Next()
	return rowToAuthor(rows)
}

func (conn *Database) Authors() ([]*Author, error) {
	rows, err := conn.db.Query("select handle, email, type, status from author")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	authors := make([]*Author, 0)

	for rows.Next() {
		author, err := rowToAuthor(rows)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func (conn *Database) CreateAuthor(handle, email, password string) error {
	_, err := conn.db.Exec(
		"insert into author (handle, email, password) values ($1, $2, $3)",
		handle, email, password,
	)

	return err
}

func rowToAuthor(rows *sql.Rows) (*Author, error) {
	var a Author
	err := rows.Scan(&a.Handle, &a.Email, &a.Status, &a.Type)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
