package main

import (
	"flag"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	DSN           string `yaml:"dsn"`
	ListenAddress string `yaml:"listen-address"`
}

func mustParseConfig() config {
	var (
		file string
		c    config
	)

	flag.StringVar(&file, "config", "", "Config file")
	flag.Parse()

	if file == "" {
		log.Warn("Config file not found in parameters, fallback to default: config.yaml")
		file = "config.yaml"
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if err = yaml.Unmarshal(bytes, &c); err != nil {
		log.Fatal(err)
	}
	log.Infof("Loaded config %#v", c)

	return c
}
