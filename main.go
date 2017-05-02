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
// Construction
//-----------------------------------------------------------------------------

func mkResources() *internal.Resources {
	log.Println("Constructing resources.")
	r, err := internal.NewResources()
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func mkConfig(resources *internal.Resources) *internal.AppConfig {
	log.Println("Constructing application configuration.")
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

	return c
}

func mkDatabase(config *internal.AppConfig) *internal.Database {
	log.Println("Constructing database connection.")
	d := internal.NewDatabase(config.Storage)

	d.Connect()
	return d
}

func mkGraphAPI(database *internal.Database) *internal.GraphAPI {
	log.Println("Constructing GraphQL API.")

	api, err := internal.NewApi(database)
	if err != nil {
		panic(err)
	}
	return api
}

func mkWebApp(resources *internal.Resources, database *internal.Database,
	graphapi *internal.GraphAPI) *http.ServeMux {
	log.Println("Constructing web app.")

	service := http.NewServeMux()

	home := http.HandlerFunc(homePage(database, resources))
	gql := http.HandlerFunc(graphQlClientPage(resources))

	service.Handle("/graphql", gql)
	service.Handle("/query", &relay.Handler{Schema: graphapi.Schema})
	service.Handle("/", home)

	return service
}

//-----------------------------------------------------------------------------
// Boostraap
//-----------------------------------------------------------------------------

func main() {

	log.Println("Welcome to Webl")

	resources := mkResources()
	config := mkConfig(resources)
	database := mkDatabase(config)
	graphapi := mkGraphAPI(database)
	app := mkWebApp(resources, database, graphapi)

	log.Printf("Web access -> http://localhost:%s/\n", config.Web.Port)
	log.Printf("GraphQL explorer access -> http://localhost:%s/graphql\n", config.Web.Port)
	log.Printf("Query API access -> http://localhost:%s/query\n", config.Web.Port)

	log.Fatal(http.ListenAndServe(":"+config.Web.Port, app))
}
