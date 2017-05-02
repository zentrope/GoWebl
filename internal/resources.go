package internal

import (
	"log"

	rice "github.com/GeertJohan/go.rice"
)

var resourceBox *rice.Box
var publicResources *rice.Box

func Resources() *rice.Box {
	if resourceBox == nil {
		log.Println("Resources loaded.")
		resourceBox = rice.MustFindBox("../resources")
	}
	return resourceBox
}

func PublicResources() *rice.Box {
	if publicResources == nil {
		log.Println("Public resources loaded.")
		publicResources = rice.MustFindBox("../resources/public")
	}
	return publicResources
}
