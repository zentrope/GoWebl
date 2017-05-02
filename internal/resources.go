package internal

import (
	"html/template"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

type Resources struct {
	Private *rice.Box
	Public  *rice.Box
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

	return &Resources{private, public}, nil
}

func (r *Resources) ResolveTemplate(name string) (*template.Template, error) {
	templateString, err := r.Private.String(name + ".template")

	if err != nil {
		return nil, err
	}

	return template.New(name).Parse(templateString)
}

func (r *Resources) PublicFileServer() http.Handler {
	return http.FileServer(r.Public.HTTPBox())
}
