// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
)

type HomeData struct {
	Posts []*TemplatePost
}

type ArchiveData struct {
	Entries []*TemplateArchiveEntry
}

type TemplatePost struct {
	UUID        string
	Author      string
	Email       string
	DateCreated string
	DateUpdated string
	Status      string
	Slugline    string
	Text        template.HTML
}

type TemplateArchiveEntry struct {
	UUID        string
	DateCreated string
	DateUpdated string
	Slugline    string
	Author      string
}

func xformArchiveEntry(e *ArchiveEntry) *TemplateArchiveEntry {
	return &TemplateArchiveEntry{
		e.UUID,
		e.DateCreated.Format("Jan 2, 2006"),
		e.DateUpdated.Format("Jan 2, 2006"),
		e.Slugline,
		e.Author,
	}
}

func toMarkdown(data string) string {
	return string(blackfriday.MarkdownCommon([]byte(data)))
}

func xformTemplatePost(p *LatestPost) *TemplatePost {
	return &TemplatePost{
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

func logRequest(r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL.Path)
}

func isIndexPath(prefix string, r *http.Request) bool {
	path := r.URL.Path
	return (path == prefix) || strings.HasSuffix(path, "/index.html")
}

func StaticPage(resources *Resources) http.HandlerFunc {
	fs := resources.AdminFileServer()
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		fs.ServeHTTP(w, r)
	}
}

func AdminPage(resources *Resources) http.HandlerFunc {
	fs := resources.AdminFileServer()

	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		if resources.AdminFileExists(r.URL.Path[1:]) {
			fs.ServeHTTP(w, r)
			return
		}

		page, err := resources.Admin.String("index.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, page)
	}
}

func HomePage(database *Database, resources *Resources) http.HandlerFunc {

	page, err := resources.ResolveTemplate("index.html")

	if err != nil {
		log.Fatal(err)
	}

	fs := resources.PublicFileServer()

	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		if !isIndexPath("/", r) {
			fs.ServeHTTP(w, r)
			return
		}

		posts, err := database.LatestPosts(40)

		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		rPosts := make([]*TemplatePost, 0)
		for _, p := range posts {
			rPosts = append(rPosts, xformTemplatePost(p))
		}

		data := &HomeData{rPosts}

		if err := page.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
}

func PostPage(database *Database, resources *Resources) http.HandlerFunc {
	page, err := resources.ResolveTemplate("post.html")

	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		logRequest(r)

		uuid := strings.Split(r.URL.Path, "/")[2]

		post, err := database.FocusPost(uuid)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		data := xformTemplatePost(post)

		if err := page.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
}

func ArchivePage(database *Database, resources *Resources) http.HandlerFunc {

	page, err := resources.ResolveTemplate("archive.html")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		entries, err := database.ArchiveEntries()
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		data := make([]*TemplateArchiveEntry, 0)
		for _, e := range entries {
			data = append(data, xformArchiveEntry(e))
		}

		values := &ArchiveData{data}

		if err := page.Execute(w, values); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
}

func GraphQlClientPage(resources *Resources) http.HandlerFunc {
	page, err := resources.PrivateString("graphql.html")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		fmt.Fprintf(w, page)
	}
}
