package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	DSN              string `yaml:"dsn"`
	ListenAddress    string `yaml:"listen-address"`
	SlackToken       string `yaml:"slack-token"`
	SlackChannel     string `yaml:"slack-channel"`
	TestSlackChannel string `yaml:"test-slack-channel"`
	MovioRealmID     int    `yaml:"movio-realm-id"`
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

	// OS env takes priority over config file (heroku doesn't support config files)
	// heroku defines an ENV var that is the port that should be exposed
	port := os.Getenv("PORT")
	// prod dsn is set in heroku via config var
	dsn := os.Getenv("HEROKU_POSTGRESQL_PINK_URL")

	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	testSlackChannel := os.Getenv("TEST_SLACK_CHANNEL")
	movioRealmID := os.Getenv("MOVIO_REALM_ID")

	if dsn != "" {
		c.DSN = dsn
	}

	if port != "" {
		c.ListenAddress = ":" + port
	}

	if slackToken != "" {
		c.SlackToken = slackToken
	}

	if slackChannel != "" {
		c.SlackChannel = slackChannel
	}

	if testSlackChannel != "" {
		c.TestSlackChannel = testSlackChannel
	}

	if movioRealmID != "" {
		c.MovioRealmID, _ = strconv.Atoi(movioRealmID)

	}

	log.Infof("Loaded config %#v", c)

	return c
}
