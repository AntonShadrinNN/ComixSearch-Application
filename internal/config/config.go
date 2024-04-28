// Package config cotains configuration settings for the application.
package config

import (
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	UrlArchive string `yaml:"urlArchive"`
	UrlComic   string `yaml:"urlComic"`
	DbConn     string `yaml:"dbConn"`
	Port       string `yaml:"port"`
	MaxProc    int    `yaml:"maxProc"`
}

func parseConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return &Config{}, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var conf Config
	err = decoder.Decode(&conf)
	if err != nil {
		return &Config{}, err
	}

	return &conf, err
}

func GetConfig(fileName string) (Config, error) {
	cnf, err := parseConfig(fileName)
	if err != nil {
		return Config{}, err
	}

	cnf.MaxProc = runtime.NumCPU()
	return *cnf, nil
}
