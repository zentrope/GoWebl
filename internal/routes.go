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
	"os"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

type WebApplication struct {
	router    *http.ServeMux
	resources *Resources
	graphql   *GraphAPI
	database  *Database
}

type ResourceKey string

const (
	DB_KEY   = ResourceKey("db")
	API_KEY  = ResourceKey("api")
	SITE_KEY = ResourceKey("conf")
	RES_KEY  = ResourceKey("res")
	AUTH_KEY = ResourceKey("auth")
)

func (app *WebApplication) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// log.Printf("> %v %v", r.Method, r.URL.String())

	w.Header().Set("Access-Control-Allow-Origin", getOriginDomain(r))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	site, err := app.database.GetSiteConfig()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx1 := context.WithValue(r.Context(), DB_KEY, app.database)
	ctx2 := context.WithValue(ctx1, API_KEY, app.graphql)
	ctx3 := context.WithValue(ctx2, SITE_KEY, site)
	ctx4 := context.WithValue(ctx3, RES_KEY, app.resources)

	app.router.ServeHTTP(w, r.WithContext(ctx4))
}

func NewWebApplication(config *AppConfig, resources *Resources,
	database *Database, graphapi *GraphAPI) *WebApplication {

	service := http.NewServeMux()

	// GraphQL
	service.HandleFunc("/graphql", graphQlClientPage)
	service.HandleFunc("/static/", staticPage)
	service.HandleFunc("/vendor/", staticPage)

	// Admin Post Manager
	service.HandleFunc("/admin/", adminPage)

	// Public Blog Routes
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
	}
}

type homeData struct {
	Posts []*TemplatePost
	Site  *SiteConfig
}

type archiveData struct {
	Entries []*TemplateArchiveEntry
	Site    *SiteConfig
}

type postData struct {
	Post *TemplatePost
	Site *SiteConfig
}

type TemplatePost struct {
	UUID        string
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
}

func westCoastTZ(date time.Time) time.Time {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Println("Error", err)
		return date
	}
	return date.In(loc)
}

func xformArchiveEntry(e *ArchiveEntry) *TemplateArchiveEntry {
	fmt := "January 2, 2006"
	return &TemplateArchiveEntry{
		e.UUID,
		westCoastTZ(e.DateCreated).Format(fmt),
		westCoastTZ(e.DateUpdated).Format(fmt),
		e.Slugline,
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
		westCoastTZ(p.DateCreated).Format("January 2, 2006"),
		westCoastTZ(p.DateUpdated).Format("January 2, 2006"),
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
	site := resolveSite(r)
	database := resolveDb(r)

	posts, err := database.LatestPosts(40)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	feed, err := NewRSSFeed(site, posts)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/xml")
	fmt.Fprintf(w, feed)
}

func jsonFeed(w http.ResponseWriter, r *http.Request) {
	site := resolveSite(r)
	database := resolveDb(r)

	posts, err := database.LatestPosts(40)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	feed, err := NewJSONFeed(site, posts)
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
	site := resolveSite(r)

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

	data := &homeData{Site: site, Posts: rPosts}

	if err := page.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {

	resources := resolveRes(r)
	database := resolveDb(r)
	site := resolveSite(r)

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

	data := &postData{Post: xformTemplatePost(post), Site: site}

	if err := page.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func archivePage(w http.ResponseWriter, r *http.Request) {

	resources := resolveRes(r)
	database := resolveDb(r)
	site := resolveSite(r)

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

	values := &archiveData{Entries: data, Site: site}

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

func getOriginDomain(r *http.Request) string {

	pattern := "http://%v:3000"

	urlhost := strings.ToLower(r.Host)
	localhost := fmt.Sprintf(pattern, "localhost")
	if strings.Contains(urlhost, "localhost") {
		return localhost
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to get hostname: ", err)
		return localhost
	}

	return strings.ToLower(fmt.Sprintf(pattern, hostname))
}

func resolveDb(r *http.Request) *Database {
	return r.Context().Value(DB_KEY).(*Database)
}

func resolveApi(r *http.Request) *GraphAPI {
	return r.Context().Value(API_KEY).(*GraphAPI)
}

func resolveSite(r *http.Request) *SiteConfig {
	return r.Context().Value(SITE_KEY).(*SiteConfig)
}

func resolveRes(r *http.Request) *Resources {
	return r.Context().Value(RES_KEY).(*Resources)
}
