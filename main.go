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
	"syscall"

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
	d.MustConnect()
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

	admin := http.HandlerFunc(internal.AdminPage(resources))
	api := http.HandlerFunc(internal.QueryAPI(graphapi))
	archive := http.HandlerFunc(internal.ArchivePage(database, resources))
	gql := http.HandlerFunc(internal.GraphQlClientPage(resources))
	home := http.HandlerFunc(internal.HomePage(database, resources))
	post := http.HandlerFunc(internal.PostPage(database, resources))
	static := http.HandlerFunc(internal.StaticPage(resources))

	// GraphQL
	service.Handle("/graphql", gql)
	service.Handle("/static/", static)
	service.Handle("/vendor/", static)

	// Admin Post Manager
	service.Handle("/admin/", admin)

	// public blog routes
	service.Handle("/archive", archive)
	service.Handle("/query", api)
	service.Handle("/post/", post)
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
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	v := <-sigChan
	log.Printf("Signal: %v\n", v)
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
		log.Println("Shutdown")
		log.Println("- Server shutdown.")
		server.Close()
		log.Println("- Database disconnect.")
		database.Disconnect()
	})

	log.Println("System halt.")
}
