package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	SMPP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"smpp"`
	REST struct {
		Url   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"rest"`
	SYSTEM struct {
		Log string `yaml:"log"`
	} `yaml:"system"`
}

func readConfigFile(cfg *Config) {
	f, err := os.Open("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
