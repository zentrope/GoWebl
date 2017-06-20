// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"database/sql"
	"encoding/hex"
	"errors"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Author struct {
	Uuid   string
	Name   string
	Email  string
	Type   string
	Status string
}

func (conn *Database) Authentic(email, password string) (*Author, error) {
	const query = "select uuid, password from author where lower(email)=lower($1)"
	rows, err := conn.db.Query(query, email)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("User not found.")
	}

	var hash string
	var authorUuid string
	err = rows.Scan(&authorUuid, &hash)
	if err != nil {
		return nil, err
	}

	decoded, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(decoded, []byte(password))
	if err != nil {
		return nil, err
	}

	return conn.Author(authorUuid)
}

func (conn *Database) Author(authorUuid string) (*Author, error) {
	const query = "select uuid, name, email, type, status from author where uuid=$1"
	rows, err := conn.db.Query(query, authorUuid)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	rows.Next()
	return rowToAuthor(rows)
}

func (conn *Database) Authors() ([]*Author, error) {
	rows, err := conn.db.Query("select uuid, name, email, type, status from author")
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

func rowToAuthor(rows *sql.Rows) (*Author, error) {
	var a Author
	err := rows.Scan(&a.Uuid, &a.Name, &a.Email, &a.Status, &a.Type)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
