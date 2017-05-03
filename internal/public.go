// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/russross/blackfriday"
)

type HomeData struct {
	Posts []*HtmlPost
}

type HtmlPost struct {
	UUID        string
	Author      string
	Email       string
	DateCreated string
	DateUpdated string
	Status      string
	Slugline    string
	Text        template.HTML
}

func toMarkdown(data string) string {
	return string(blackfriday.MarkdownBasic([]byte(data)))
}

func templatize(p *LatestPost) *HtmlPost {
	return &HtmlPost{
		p.UUID,
		p.Author,
		p.Email,
		p.DateCreated.Format("Jan 2, 2006"),
		p.DateUpdated.Format("Jan 2, 2006"),
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

		posts, err := database.LatestPosts(40)

		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		rPosts := make([]*HtmlPost, 0)
		for _, p := range posts {
			rPosts = append(rPosts, templatize(p))
		}

		data := &HomeData{rPosts}

		if err := page.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
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
