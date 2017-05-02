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

	rice "github.com/GeertJohan/go.rice"
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

func resolveTemplate(name string) (*template.Template, error) {
	resources := rice.MustFindBox("resources")

	templateString, err := resources.String(name + ".template")

	if err != nil {
		return nil, err
	}

	return template.New(name).Parse(templateString)
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

// func MkHandler(location string) http.HandlerFunc {

//	// This craziness allows me to serve index.html when routes
//	// are otherwise not found.

//	box := rice.MustFindBox(location)
//	fs := http.FileServer(box.HTTPBox())

//	return func(w http.ResponseWriter, r *http.Request) {
//		content := r.URL.Path[1:]

//		f, err := box.Open(content)

//		if err == nil {
//			f.Close()
//			fs.ServeHTTP(w, r)
//			return
//		}

//		body, _ := box.Bytes("index.html")
//		w.Header().Set("Content-Type", "text/html")
//		w.WriteHeader(http.StatusOK)
//		w.Write(body)
//	}
// }

func homePage(database *internal.Database) http.HandlerFunc {

	page, err := resolveTemplate("index.html")

	if err != nil {
		log.Fatal(err)
	}

	box := rice.MustFindBox("resources/public")
	fs := http.FileServer(box.HTTPBox())

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

func graphQlClientPage() http.HandlerFunc {
	page, err := resolveTemplate("graphql.html")
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
	log.Println("Initializing application configuration.")
	var overrideConfigFile string
	flag.StringVar(&overrideConfigFile, "c", internal.DefaultConfigFile,
		"Path to configuration override file.")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	c, err := internal.LoadConfigFile(overrideConfigFile)
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

func main() {
	defer database.Disconnect()

	home := http.HandlerFunc(homePage(database))
	gql := http.HandlerFunc(graphQlClientPage())

	http.Handle("/graphql", gql)
	http.Handle("/query", &relay.Handler{Schema: graphapi.Schema})
	http.Handle("/", home)

	fmt.Printf("Hello Webl\n")
	fmt.Printf(" - web     -> http://localhost:%s/\n", config.Web.Port)
	fmt.Printf(" - graphql -> http://localhost:%s/graphql\n", config.Web.Port)
	fmt.Printf(" - query   -> http://localhost:%s/query\n", config.Web.Port)

	log.Fatal(http.ListenAndServe(":"+config.Web.Port, nil))
}
