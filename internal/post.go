// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"database/sql"
	"fmt"
	"time"
)

//-----------------------------------------------------------------------------
// Queries for Public Pages
//-----------------------------------------------------------------------------

type LatestPost struct {
	UUID        string
	DateCreated time.Time
	DateUpdated time.Time
	Status      string
	Slugline    string
	Author      string
	Email       string
	Text        string
}

func (conn *Database) FocusPost(uuid string) (*LatestPost, error) {
	var query = `
	 select
		 p.uuid, p.date_created, p.date_updated, p.status,
		 p.slugline, a.handle as author, a.email, p.text
	 from
		 post p, author a
	 where
		 p.uuid = $1
		 and p.author=a.handle
		 and p.status='published'
	`

	rows, err := conn.db.Query(query, uuid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rows.Next()
	var p LatestPost

	err = rows.Scan(
		&p.UUID,
		&p.DateCreated,
		&p.DateUpdated,
		&p.Status,
		&p.Slugline,
		&p.Author,
		&p.Email,
		&p.Text,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (conn *Database) LatestPosts(limit int) ([]*LatestPost, error) {
	var query = `
	 select
		 p.uuid, p.date_created, p.date_updated, p.status,
		 p.slugline, a.handle as author, a.email, p.text
	 from
		 post p, author a
	 where
		 p.author=a.handle
		 and p.status='published'
	 order
		 by date_created desc
	 limit $1`

	rows, err := conn.db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := make([]*LatestPost, 0)

	for rows.Next() {
		var p LatestPost
		err := rows.Scan(
			&p.UUID,
			&p.DateCreated,
			&p.DateUpdated,
			&p.Status,
			&p.Slugline,
			&p.Author,
			&p.Email,
			&p.Text,
		)

		if err != nil {
			return make([]*LatestPost, 0), err
		}

		posts = append(posts, &p)
	}

	return posts, nil
}

type ArchiveEntry struct {
	UUID        string
	DateCreated time.Time
	DateUpdated time.Time
	Slugline    string
	Author      string
}

func (conn *Database) ArchiveEntries() ([]*ArchiveEntry, error) {

	var query = `
		select
			uuid, date_created, date_updated, slugline, author
		from
			post
		where
			status='published'
		order by
			date_created desc;
	`

	rows, err := conn.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	entries := make([]*ArchiveEntry, 0)

	for rows.Next() {
		var e ArchiveEntry
		err := rows.Scan(
			&e.UUID,
			&e.DateCreated,
			&e.DateUpdated,
			&e.Slugline,
			&e.Author,
		)

		if err != nil {
			return make([]*ArchiveEntry, 0), err
		}

		entries = append(entries, &e)
	}

	return entries, nil

}

//-----------------------------------------------------------------------------
// Queries for GraphQL
//-----------------------------------------------------------------------------

type Post struct {
	UUID        string
	Author      string
	DateCreated time.Time
	DateUpdated time.Time
	Status      string
	Slugline    string
	Text        string
}

func (conn *Database) Posts() ([]*Post, error) {
	q := mkPostSql("")
	return conn.postQuery(q)
}

func (conn *Database) PostsByAuthor(handle string) ([]*Post, error) {
	return conn.postQuery(
		mkPostSql("where lower(author)=lower($1)"),
		handle,
	)
}

func (conn *Database) postQuery(query string, args ...interface{}) ([]*Post, error) {

	rows, err := conn.db.Query(query, args...)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	results, err := rowsToPosts(rows)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func mkPostSql(where string) string {
	q := "select uuid, author, date_created, date_updated, status, slugline, text from post %s"
	return fmt.Sprintf(q, where)
}

func rowsToPosts(rows *sql.Rows) ([]*Post, error) {

	posts := make([]*Post, 0)

	for rows.Next() {
		post, err := rowToPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func rowToPost(rows *sql.Rows) (*Post, error) {
	var p Post
	err := rows.Scan(
		&p.UUID,
		&p.Author,
		&p.DateCreated,
		&p.DateUpdated,
		&p.Status,
		&p.Slugline,
		&p.Text,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
