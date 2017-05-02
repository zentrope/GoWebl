package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/neelance/graphql-go/relay"
	"github.com/russross/blackfriday"
	"github.com/zentrope/webl/internal"
)

type HomeData struct {
	Authors []*internal.Author
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

func templatize(p *internal.Post) *HtmlPost {
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

func homePage(database *internal.Database, resources *internal.Resources) http.HandlerFunc {

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

func graphQlClientPage(resources *internal.Resources) http.HandlerFunc {
	page, err := resources.ResolveTemplate("graphql.html")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := make([]interface{}, 0)
		page.Execute(w, &data)
	}
}

//-----------------------------------------------------------------------------
// Application Components and Resources
//-----------------------------------------------------------------------------

var resources *internal.Resources
var config *internal.AppConfig
var database *internal.Database
var graphapi *internal.GraphAPI

//-----------------------------------------------------------------------------
// Initializers
//-----------------------------------------------------------------------------

func init() {
	log.Println("Welcome to Webl")
}

func init() {
	log.Println("Initializing resources.")
	r, err := internal.NewResources()
	if err != nil {
		log.Fatal(err)
	}
	resources = r
}

func init() {
	log.Println("Initializing application configuration.")
	var overrideConfigFile string
	flag.StringVar(&overrideConfigFile, "c", internal.DefaultConfigFile,
		"Path to configuration override file.")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	c, err := internal.LoadConfigFile(overrideConfigFile, resources)
	if err != nil {
		panic(err)
	}

	config = c
}

func init() {
	log.Println("Initializing database connection.")
	d := internal.NewDatabase(config.Storage)

	d.Connect()

	database = d
}

func init() {
	log.Println("Intializing GraphQL API")

	a, err := internal.NewApi(database)
	if err != nil {
		panic(err)
	}
	graphapi = a
}

//-----------------------------------------------------------------------------
// Bootstrap
//-----------------------------------------------------------------------------

func main() {

	home := http.HandlerFunc(homePage(database, resources))
	gql := http.HandlerFunc(graphQlClientPage(resources))

	http.Handle("/graphql", gql)
	http.Handle("/query", &relay.Handler{Schema: graphapi.Schema})
	http.Handle("/", home)

	fmt.Printf("Hello Webl\n")
	fmt.Printf(" - web     -> http://localhost:%s/\n", config.Web.Port)
	fmt.Printf(" - graphql -> http://localhost:%s/graphql\n", config.Web.Port)
	fmt.Printf(" - query   -> http://localhost:%s/query\n", config.Web.Port)

	log.Fatal(http.ListenAndServe(":"+config.Web.Port, nil))
}
