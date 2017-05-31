// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"encoding/json"
	"io/ioutil"

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
	Port    string `json:"port,omitempty"`
	Title   string `json:"title,omitempty"`
	BaseURL string `json:"base-url,omitempty"`
}

type AppConfig struct {
	Storage StorageConfig `json:"storage,omitempty"`
	Web     WebConfig     `json:"web,omitempty"`
}

const DefaultConfigFile = "resources/config.json"

func LoadConfigFile(pathToOverride string, resources *Resources) (*AppConfig, error) {

	config, err := loadDefaultConfigFile(resources)
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
		return nil, err
	}

	return &override, nil
}

//-----------------------------------------------------------------------------
// Implementation
//-----------------------------------------------------------------------------

func loadDefaultConfigFile(resources *Resources) (AppConfig, error) {
	contents, err := resources.Private.String("config.json")

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
