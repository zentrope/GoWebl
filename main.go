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
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/zentrope/webl/server"
)

//-----------------------------------------------------------------------------
// Construction
//-----------------------------------------------------------------------------

var resourceDir, adminDir, overrideFile, assetDir string

func init() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	resourcePath := filepath.Join(dir, "resources")
	assetPath := filepath.Join(dir, "assets")
	adminPath := filepath.Join(dir, "admin")

	flag.StringVar(&resourceDir, "resources", resourcePath, "Path to scripts, templates.")
	flag.StringVar(&assetDir, "assets", assetPath, "Path to web assets.")
	flag.StringVar(&adminDir, "app", adminPath, "Path to admin web app.")
	flag.StringVar(&overrideFile, "c", "", "Path to configuration override file.")

	flag.Parse()

	log.Println("Config:")
	t := "- %-9v '%v'"
	log.Printf(t, "resource:", resourceDir)
	log.Printf(t, "admin:", adminDir)
	log.Printf(t, "asset:", assetDir)
	log.Printf(t, "override:", overrideFile)
}

func mkResources() server.Resources {
	log.Println("Constructing resources.")

	r, err := server.NewResources(resourceDir, assetDir, adminDir)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func mkConfig(resources server.Resources) *server.AppConfig {
	log.Println("Constructing application configuration.")

	c, err := server.LoadConfigFile(overrideFile)
	if err != nil {
		panic(err)
	}

	return c
}

func mkDatabase(config *server.AppConfig) *server.Database {
	log.Println("Constructing database connection.")
	d := server.NewDatabase(config.Storage)
	d.MustConnect()
	return d
}

func mkGraphAPI(database *server.Database) *server.GraphAPI {
	log.Println("Constructing GraphQL API.")

	api, err := server.NewApi(database)
	if err != nil {
		panic(err)
	}
	return api
}

func mkWebApp(
	config *server.AppConfig,
	resources server.Resources,
	database *server.Database,
	graphapi *server.GraphAPI,
) *http.Server {

	app := server.NewWebApplication(config, resources, database, graphapi)

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

func notify(config *server.AppConfig) {
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
