package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Mods []Mod `json:"mods"`
}

type Mod struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Loader   string `json:"loader"`
	Filename string `json:"filename,omitempty"`

	Replacers map[string]string `json:"replacers"`
}

func parseConf(path string) (*Config, error) {
	var conf Config

	f, err := os.ReadFile(path)
	if err != nil {
		return &conf, err
	}

	err = json.Unmarshal(f, &conf)
	return &conf, err
}
