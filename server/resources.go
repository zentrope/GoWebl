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
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

type Resources struct {
	Private *rice.Box
	Public  *rice.Box
	Admin   *rice.Box
}

func NewResources() (*Resources, error) {

	private, err := rice.FindBox("../resources")
	if err != nil {
		return nil, err
	}

	public, err := rice.FindBox("../resources/public")
	if err != nil {
		return nil, err
	}

	admin, err := rice.FindBox("../admin/build")
	if err != nil {
		return nil, err
	}

	return &Resources{private, public, admin}, nil
}

var cache = NewCache()

func (r *Resources) ResolveTemplate(name string) (*template.Template, error) {

	t, err := cache.GetOrSet(name, func() (interface{}, error) {
		templateString, err := r.Private.String(name)
		if err != nil {
			return nil, err
		}
		return template.New(name).Parse(templateString)
	})

	if err != nil {
		return nil, err
	}

	return t.(*template.Template), nil
}

func (r *Resources) AdminFileExists(name string) bool {
	_, err := r.Admin.Open(name)
	return err == nil
}

func (r *Resources) PublicFileServer() http.Handler {
	return http.FileServer(r.Public.HTTPBox())
}

func (r *Resources) AdminFileServer() http.Handler {
	return http.FileServer(r.Admin.HTTPBox())
}

func (r *Resources) PrivateString(name string) (string, error) {
	return r.Private.String(name)
}
