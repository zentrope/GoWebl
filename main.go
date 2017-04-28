package main

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/neelance/graphql-go/relay"
	"github.com/russross/blackfriday"
	"github.com/zentrope/webl/api"
	"github.com/zentrope/webl/database"
	"github.com/zentrope/webl/resources"
)

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

func homePage(database *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authors := database.Authors()
		posts := database.Posts()

		rPosts := make([]*HtmlPost, 0)
		for _, p := range posts {
			rPosts = append(rPosts, templatize(p))
		}

		data := &HomeData{authors, rPosts}

		page := resources.HomePageTemplate
		page.Execute(w, data)
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

	http.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(resources.GraphQlPage)
	}))

	http.Handle("/query", &relay.Handler{Schema: api})

	http.Handle("/index.html", home)
	http.Handle("/", home)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
