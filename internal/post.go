package internal

import (
	"database/sql"
	"log"
	"time"
)

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
