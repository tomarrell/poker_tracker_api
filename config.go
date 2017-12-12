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
		flag.Usage()
		log.Fatal("Config file is required")
	}

	byts, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if err = yaml.Unmarshal(byts, &c); err != nil {
		log.Fatal(err)
	}
	log.Infof("Loaded config %#v", c)

	return c
}
