// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

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
