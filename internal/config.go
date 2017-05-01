package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rice "github.com/GeertJohan/go.rice"
	"github.com/imdario/mergo"
)

//-----------------------------------------------------------------------------
// Exports
//-----------------------------------------------------------------------------

type StorageConfig struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
}

type WebConfig struct {
	Port string `json:"port,omitempty"`
}

type AppConfig struct {
	Storage StorageConfig `json:"storage,omitempty"`
	Web     WebConfig     `json:"web,omitempty"`
}

const DefaultConfigFile = "resources/config.json"

func LoadConfigFile(pathToOverride string) (*AppConfig, error) {

	config, err := loadDefaultConfigFile()
	if err != nil {
		return nil, err
	}

	if pathToOverride == DefaultConfigFile {
		return &config, nil
	}

	override, err := loadConfigFile(pathToOverride)
	if err != nil {
		return nil, err
	}

	if err := mergo.Merge(&override, config); err != nil {
		fmt.Printf("Merge Error: %v\n", err)
	}

	return &override, nil
}

//-----------------------------------------------------------------------------
// Implementation
//-----------------------------------------------------------------------------

var resourceBox *rice.Box

func resources() *rice.Box {
	if resourceBox == nil {
		fmt.Println("resources loaded")
		resourceBox = rice.MustFindBox("../resources")
	}
	return resourceBox
}

func loadDefaultConfigFile() (AppConfig, error) {
	contents, err := resources().String("config.json")

	if err != nil {
		return AppConfig{}, err
	}

	var config AppConfig
	if err := json.Unmarshal([]byte(contents), &config); err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

func loadConfigFile(path string) (AppConfig, error) {

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return AppConfig{}, err
	}

	var config AppConfig
	if err := json.Unmarshal(contents, &config); err != nil {
		return AppConfig{}, err
	}

	return config, nil
}
