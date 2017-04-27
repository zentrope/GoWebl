package main

import (
	"fmt"
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

func MkHandler(location string) http.HandlerFunc {

	// This craziness allows me to serve index.html when routes
	// are otherwise not found.

	box := rice.MustFindBox(location)
	fs := http.FileServer(box.HTTPBox())

	return func(w http.ResponseWriter, r *http.Request) {
		content := r.URL.Path[1:]

		f, err := box.Open(content)

		if err == nil {
			f.Close()
			fs.ServeHTTP(w, r)
			return
		}

		body, _ := box.Bytes("index.html")
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func main() {
	fmt.Println("Hello Server")

	location = "../client"
	fmt.Println("location:", location)

	http.Handle("/", MkHandler(location))

	http.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("graphql")
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<h1>Not available</h1>")
	}))
	log.Fatal(http.ListenAndServe(":12345", nil))
}
