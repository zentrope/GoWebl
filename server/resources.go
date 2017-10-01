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

package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Resources represent static files, templates, images, etc.
type Resources struct {
	privateDir string
	publicDir  string
	adminDir   string
}

// NewResources returns a new instance for the resource
// manager for loading static files.
func NewResources() (Resources, error) {
	log.Printf("WARNING: resources are using hard coded paths.")
	pr, _ := filepath.Abs("./resources")
	pu, _ := filepath.Abs("./resources/public")
	ad, _ := filepath.Abs("./admin/build")

	return Resources{
		privateDir: pr,
		publicDir:  pu,
		adminDir:   ad,
	}, nil
}

func (r Resources) resolveTemplate(name string) (*template.Template, error) {
	templateString, err := r.privateString(name)
	if err != nil {
		return &template.Template{}, err
	}
	return template.New(name).Parse(templateString)
}

func (r Resources) adminFileExists(name string) bool {
	return fileExists(r.adminDir, name)
}

func (r Resources) publicFileServer() http.Handler {
	return http.FileServer(http.Dir(r.publicDir))
}

func (r Resources) adminFileServer() http.Handler {
	return http.FileServer(http.Dir(r.adminDir))
}

func (r Resources) adminString(name string) (string, error) {
	return loadFile(r.adminDir, name)
}

func (r Resources) privateString(name string) (string, error) {
	return loadFile(r.privateDir, name)
}

func loadFile(root, name string) (string, error) {
	path := filepath.Join(root, name)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func fileExists(root, name string) bool {
	path := filepath.Join(root, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
