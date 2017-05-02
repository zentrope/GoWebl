// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/russross/blackfriday"
)

type HomeData struct {
	Authors []*Author
	Posts   []*HtmlPost
}

type HtmlPost struct {
	Id          string
	Author      string
	DateCreated string
	DateUpdated string
	Status      string
	Slugline    string
	Text        template.HTML
}

func idToStr(id int) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.Itoa(id))))
}

func toMarkdown(data string) string {
	return string(blackfriday.MarkdownBasic([]byte(data)))
}

func templatize(p *Post) *HtmlPost {
	return &HtmlPost{
		idToStr(p.Id),
		p.Author,
		p.DateCreated.Format("Jan 2, 2006"),
		p.DateUpdated.Format(time.RFC822),
		p.Status,
		p.Slugline,
		template.HTML(toMarkdown(p.Text)),
	}
}

func HomePage(database *Database, resources *Resources) http.HandlerFunc {

	page, err := resources.ResolveTemplate("index.html")

	if err != nil {
		log.Fatal(err)
	}

	fs := resources.PublicFileServer()

	return func(w http.ResponseWriter, r *http.Request) {

		if (r.URL.Path != "/") && (r.URL.Path != "/index.html") {
			fs.ServeHTTP(w, r)
			return
		}

		authors := database.Authors()
		posts := database.Posts()

		rPosts := make([]*HtmlPost, 0)
		for _, p := range posts {
			rPosts = append(rPosts, templatize(p))
		}

		data := &HomeData{authors, rPosts}

		page.Execute(w, data)
	}
}

func GraphQlClientPage(resources *Resources) http.HandlerFunc {
	page, err := resources.ResolveTemplate("graphql.html")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := make([]interface{}, 0)
		page.Execute(w, &data)
	}
}
