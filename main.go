//
// Copyright (c) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

	listen := config.Web.Listen
	if listen == "" {
		listen = "127.0.0.1"
	}

	addr := listen + ":" + config.Web.Port

	return &http.Server{
		Addr:    addr,
		Handler: app,
	}
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
	log.Printf("Ready on port %v.", config.Web.Port)
}

func main() {

	log.Println("Welcome to Webl")

	resources := mkResources()
	config := mkConfig(resources)
	database := mkDatabase(config)
	database.MustRunMigrations(resources)

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
