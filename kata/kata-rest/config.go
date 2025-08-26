package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config structures
type Request struct {
	Name       string            `yaml:"name"`
	Method     string            `yaml:"method"`
	URL        string            `yaml:"url"`
	Headers    map[string]string `yaml:"headers"`
	Parameters map[string]string `yaml:"parameters"`
	Body       string            `yaml:"body"`
}

type Collection struct {
	Name     string    `yaml:"name"`
	Requests []Request `yaml:"requests"`
}

type Config struct {
	Collections []Collection `yaml:"collections"`
}

func (app *App) loadConfig() {
	app.config = &Config{}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Error reading config.yaml: %v", err)
		// Create default empty config
		app.config.Collections = []Collection{}
		return
	}

	err = yaml.Unmarshal(data, app.config)
	if err != nil {
		log.Printf("Error parsing config.yaml: %v", err)
		app.config.Collections = []Collection{}
	}
}
