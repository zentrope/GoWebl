// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"database/sql"
	"log"
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
	 order by date_created desc
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

//-----------------------------------------------------------------------------
// Queries for GraphQL
//-----------------------------------------------------------------------------

type Post struct {
	Id          int
	Author      string
	DateCreated time.Time
	DateUpdated time.Time
	Status      string
	Slugline    string
	Text        string
}

func (conn *Database) PostStatus() []string {
	return conn.enums("post_status")
}

func (conn *Database) Posts() []*Post {
	return conn.postQuery(
		mkPostSql(""),
	)
}

func (conn *Database) PostsByAuthor(handle string) []*Post {
	return conn.postQuery(
		mkPostSql("where lower(author)=lower($1)"),
		handle,
	)
}

func (conn *Database) postQuery(query string, args ...interface{}) []*Post {

	rows, err := conn.db.Query(query, args...)

	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	return rowsToPosts(rows)
}

func mkPostSql(where string) string {
	return "select id, author, date_created, date_updated, status, slugline, text from post " + where
}

func rowsToPosts(rows *sql.Rows) []*Post {
	posts := make([]*Post, 0)

	for rows.Next() {
		posts = append(posts, rowToPost(rows))
	}

	return posts
}

func rowToPost(rows *sql.Rows) *Post {
	var p Post
	err := rows.Scan(
		&p.Id,
		&p.Author,
		&p.DateCreated,
		&p.DateUpdated,
		&p.Status,
		&p.Slugline,
		&p.Text,
	)
	if err != nil {
		log.Fatal(err)
	}
	return &p
}
