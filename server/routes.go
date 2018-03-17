//
// Copyright (c) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package server

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

	mdown "gopkg.in/russross/blackfriday.v2"
)

type WebApplication struct {
	router    *http.ServeMux
	resources Resources
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

	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Access-Control-Allow-Origin", getOriginDomain(r))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "HEAD" {
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

func NewWebApplication(config *AppConfig, resources Resources,
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
	UUID          string
	DateCreated   string
	DateUpdated   string
	DatePublished string
	Status        string
	Slugline      string
	Text          template.HTML
}

type TemplateArchiveEntry struct {
	UUID          string
	DateCreated   string
	DateUpdated   string
	DatePublished string
	Slugline      string
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
		westCoastTZ(e.DatePublished).Format(fmt),
		e.Slugline,
	}
}

func MarkdownToHtml(data string) string {

	extensions := mdown.NoIntraEmphasis |
		mdown.Tables |
		mdown.FencedCode |
		mdown.Autolink |
		mdown.Strikethrough |
		mdown.SpaceHeadings |
		mdown.HeadingIDs |
		mdown.BackslashLineBreak |
		mdown.DefinitionLists |
		mdown.Footnotes

	input := []byte(data)
	output := mdown.Run(input, mdown.WithExtensions(extensions))
	return strings.TrimSpace(string(output))
}

func xformTemplatePost(p *LatestPost) *TemplatePost {
	return &TemplatePost{
		p.UUID,
		westCoastTZ(p.DateCreated).Format("January 2, 2006"),
		westCoastTZ(p.DateUpdated).Format("January 2, 2006"),
		westCoastTZ(p.DatePublished).Format("January 2, 2006"),
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
	fs := resources.adminFileServer()
	fs.ServeHTTP(w, r)
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	resources := resolveRes(r)
	fs := resources.adminFileServer()

	if resources.adminFileExists(r.URL.Path[1:]) {
		fs.ServeHTTP(w, r)
		return
	}

	page, err := resources.adminString("index.html")
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

	fs := resources.publicFileServer()

	if !isIndexPath("/", r) {
		fs.ServeHTTP(w, r)
		return
	}

	page, err := resources.resolveTemplate("index.html")
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

	page, err := resources.resolveTemplate("post.html")
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

	page, err := resources.resolveTemplate("archive.html")
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
	page, err := resources.privateString("graphql.html")
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

func resolveRes(r *http.Request) Resources {
	return r.Context().Value(RES_KEY).(Resources)
}
