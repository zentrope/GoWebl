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

func mkWebApp(
	config *internal.AppConfig,
	resources *internal.Resources,
	database *internal.Database,
	graphapi *internal.GraphAPI,
) *http.Server {

	app := internal.NewWebApplication(config, resources, database, graphapi)

	return &http.Server{
		Addr:    ":" + config.Web.Port,
		Handler: app,
	}
}

func addSiteConfig(config *internal.AppConfig, database *internal.Database) *internal.AppConfig {
	newConfig, err := database.AppendSiteConfig(config)
	if err != nil {
		panic(err)
	}
	return newConfig
}

//-----------------------------------------------------------------------------
// Bootstrap
//-----------------------------------------------------------------------------

func blockUntilShutdownThenDo(fn func()) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	v := <-sigChan
	log.Printf("Signal: %v\n", v)
	fn()
}

func notify(config *internal.AppConfig) {
	log.Printf("Web access -> %s\n", config.Web.BaseURL)
}

func main() {

	log.Println("Welcome to Webl")

	resources := mkResources()
	config := mkConfig(resources)
	database := mkDatabase(config)
	database.MustRunMigrations(resources)
	config = addSiteConfig(config, database)

	graphapi := mkGraphAPI(database)
	server := mkWebApp(config, resources, database, graphapi)

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
