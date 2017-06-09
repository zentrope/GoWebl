// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
)

type WebApplication struct {
	router    *http.ServeMux
	resources *Resources
	graphql   *GraphAPI
	database  *Database
	config    WebConfig
}

type ResourceKey string

const DB_KEY = ResourceKey("db")
const API_KEY = ResourceKey("api")
const CONF_KEY = ResourceKey("conf")
const RES_KEY = ResourceKey("res")
const AUTH_KEY = ResourceKey("auth")

func (app *WebApplication) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// log.Printf("> %v %v", r.Method, r.URL.String())

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	ctx1 := context.WithValue(r.Context(), DB_KEY, app.database)
	ctx2 := context.WithValue(ctx1, API_KEY, app.graphql)
	ctx3 := context.WithValue(ctx2, CONF_KEY, app.config)
	ctx4 := context.WithValue(ctx3, RES_KEY, app.resources)

	app.router.ServeHTTP(w, r.WithContext(ctx4))
}

func NewWebApplication(config *AppConfig, resources *Resources,
	database *Database, graphapi *GraphAPI) *WebApplication {

	service := http.NewServeMux()

	webConfig := config.Web

	// GraphQL
	service.HandleFunc("/graphql", graphQlClientPage)
	service.HandleFunc("/static/", staticPage)
	service.HandleFunc("/vendor/", staticPage)

	// Admin Post Manager
	service.HandleFunc("/admin/", adminPage)

	// public blog routes

	service.HandleFunc("/feeds/json", jsonFeed)
	service.HandleFunc("/feeds/json/", jsonFeed)
	service.HandleFunc("/feeds/rss", rssFeed)
	service.HandleFunc("/feeds/rss/", rssFeed)

	service.HandleFunc("/archive", archivePage)
	service.HandleFunc("/query", queryApi)
	service.HandleFunc("/post/", postPage)
	service.HandleFunc("/", homePage)

	return &WebApplication{
		router:    service,
		resources: resources,
		graphql:   graphapi,
		database:  database,
		config:    webConfig,
	}
}

// ---

type homeData struct {
	Posts  []*TemplatePost
	Config WebConfig
}

type archiveData struct {
	Entries []*TemplateArchiveEntry
	Config  WebConfig
}

type postData struct {
	Post   *TemplatePost
	Config WebConfig
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
	fmt := "02-Jan-2006"
	return &TemplateArchiveEntry{
		e.UUID,
		e.DateCreated.Format(fmt),
		e.DateUpdated.Format(fmt),
		e.Slugline,
		e.Author,
	}
}

func MarkdownToHtml(data string) string {
	htmlFlags := blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	extensions := blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_FOOTNOTES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	input := []byte(data)
	options := blackfriday.Options{Extensions: extensions}

	output := blackfriday.MarkdownOptions(input, renderer, options)
	return strings.TrimSpace(string(output))
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
		template.HTML(MarkdownToHtml(p.Text)),
	}
}

func isIndexPath(prefix string, r *http.Request) bool {
	path := r.URL.Path
	return (path == prefix) || strings.HasSuffix(path, "/index.html")
}

func queryApi(w http.ResponseWriter, r *http.Request) {

	api := resolveApi(r)

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		authToken = "NO_AUTH_TOKEN"
	} else {
		authToken = strings.Replace(authToken, "Bearer ", "", 1)
	}

	authCtx := context.WithValue(r.Context(), AUTH_KEY, authToken)

	response := api.Schema.Exec(authCtx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

func staticPage(w http.ResponseWriter, r *http.Request) {
	resources := resolveRes(r)
	fs := resources.AdminFileServer()
	fs.ServeHTTP(w, r)
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	resources := resolveRes(r)
	fs := resources.AdminFileServer()

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

func rssFeed(w http.ResponseWriter, r *http.Request) {
	config := resolveConf(r)
	database := resolveDb(r)

	posts, err := database.LatestPosts(40)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	feed, err := NewRSSFeed(config, posts)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/xml")
	fmt.Fprintf(w, feed)
}

func jsonFeed(w http.ResponseWriter, r *http.Request) {
	config := resolveConf(r)
	database := resolveDb(r)

	posts, err := database.LatestPosts(40)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	feed, err := NewJSONFeed(config, posts)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, feed)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	resources := resolveRes(r)
	database := resolveDb(r)
	config := resolveConf(r)

	fs := resources.PublicFileServer()

	if !isIndexPath("/", r) {
		fs.ServeHTTP(w, r)
		return
	}

	page, err := resources.ResolveTemplate("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
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

	data := &homeData{Config: config, Posts: rPosts}

	if err := page.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {

	resources := resolveRes(r)
	database := resolveDb(r)
	config := resolveConf(r)

	page, err := resources.ResolveTemplate("post.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	uuid := strings.Split(r.URL.Path, "/")[2]

	post, err := database.FocusPost(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	data := &postData{Post: xformTemplatePost(post), Config: config}

	if err := page.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func archivePage(w http.ResponseWriter, r *http.Request) {

	resources := resolveRes(r)
	database := resolveDb(r)
	config := resolveConf(r)

	page, err := resources.ResolveTemplate("archive.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	entries, err := database.ArchiveEntries()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	data := make([]*TemplateArchiveEntry, 0)
	for _, e := range entries {
		data = append(data, xformArchiveEntry(e))
	}

	values := &archiveData{Entries: data, Config: config}

	if err := page.Execute(w, values); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func graphQlClientPage(w http.ResponseWriter, r *http.Request) {
	resources := resolveRes(r)
	page, err := resources.PrivateString("graphql.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, page)
}

//----

func resolveDb(r *http.Request) *Database {
	return r.Context().Value(DB_KEY).(*Database)
}

func resolveApi(r *http.Request) *GraphAPI {
	return r.Context().Value(API_KEY).(*GraphAPI)
}

func resolveConf(r *http.Request) WebConfig {
	return r.Context().Value(CONF_KEY).(WebConfig)
}

func resolveRes(r *http.Request) *Resources {
	return r.Context().Value(RES_KEY).(*Resources)
}
