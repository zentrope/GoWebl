package main

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/neelance/graphql-go/relay"
	"github.com/russross/blackfriday"
	"github.com/zentrope/webl/api"
	"github.com/zentrope/webl/database"
)

//	rice "github.com/GeertJohan/go.rice"

type HomeData struct {
	Authors []*database.Author
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

func templatize(p *database.Post) *HtmlPost {
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

func homePage(database *database.Database) http.HandlerFunc {

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

func main() {
	fmt.Println("Hello GraphQL")
	fmt.Println(" - web     -> http://localhost:8080/")
	fmt.Println(" - graphql -> http://localhost:8080/graphql")
	fmt.Println(" - query   -> http://localhost:8080/query")

	database := database.NewDatabase("blogsvc", "wanheda", "blogdb")

	database.Connect()
	defer database.Disconnect()

	api, err := api.NewApi(database)
	if err != nil {
		panic(err)
	}

	home := http.HandlerFunc(homePage(database))
	gql := http.HandlerFunc(graphQlClientPage())

	http.Handle("/graphql", gql)
	http.Handle("/query", &relay.Handler{Schema: api})
	http.Handle("/", home)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
