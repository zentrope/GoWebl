// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/neelance/graphql-go/relay"
	"github.com/zentrope/webl/internal"
)

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

	home := http.HandlerFunc(internal.HomePage(database, resources))
	gql := http.HandlerFunc(internal.GraphQlClientPage(resources))

	service.Handle("/graphql", gql)
	service.Handle("/query", &relay.Handler{Schema: graphapi.Schema})
	service.Handle("/", home)

	return service
}

func mkServer(config *internal.AppConfig, app *http.ServeMux) *http.Server {
	server := &http.Server{}
	server.Addr = ":" + config.Web.Port
	server.Handler = app
	return server
}

//-----------------------------------------------------------------------------
// Boostraap
//-----------------------------------------------------------------------------

func blockUntilShutdownThenDo(fn func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	fn()
}

func notify(config *internal.AppConfig) {
	p := config.Web.Port
	log.Printf("Web access -> http://localhost:%s/\n", p)
	log.Printf("GraphQL explorer access -> http://localhost:%s/graphql\n", p)
	log.Printf("Query API access -> http://localhost:%s/query\n", p)
}

func main() {

	log.Println("Welcome to Webl")

	resources := mkResources()
	config := mkConfig(resources)
	database := mkDatabase(config)
	database.MustRunMigrations(resources)

	graphapi := mkGraphAPI(database)
	app := mkWebApp(resources, database, graphapi)
	server := mkServer(config, app)

	go server.ListenAndServe()

	notify(config)

	blockUntilShutdownThenDo(func() {
		log.Println("^C")
		log.Println("Shutdown")
		log.Println("- Server shutdown.")
		server.Close()
		log.Println("- Database disconnect.")
		database.Disconnect()
	})

	log.Println("System halt.")
}
