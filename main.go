package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neelance/graphql-go/relay"
	"github.com/zentrope/webl/api"
	"github.com/zentrope/webl/database"
	"github.com/zentrope/webl/resources"
)

type HomeData struct {
	Authors []*database.Author
	Posts   []*database.Post
}

func homePage(database *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authors := database.Authors()
		posts := database.Posts()

		data := &HomeData{authors, posts}

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
