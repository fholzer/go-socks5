package main

import (
	"fmt"

	ucfg "github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/sirupsen/logrus"
)

type forwarderConfig struct {
	Type    string
	Address string
}

type ruleConfig struct {
	Subnets   []string
	Forwarder forwarderConfig
}

type rawConfiguration struct {
	Loglevel         string
	Logformat        string
	Bind             string
	Rules            []ruleConfig
	DefaultForwarder *forwarderConfig `config:"defaultForwarder"`
}

type Configuration struct {
	Loglevel         logrus.Level
	Bind             string
	Rules            []Rule
	DefaultForwarder *Forwarder
}

var (
	defaultConfig = rawConfiguration{
		Loglevel:         "info",
		Logformat:        "text",
		Bind:             "127.0.0.1:5757",
		Rules:            nil,
		DefaultForwarder: nil,
	}
)

func ParseConfig(filename string) (*Configuration, error) {
	config, err := yaml.NewConfigWithFile(filename, ucfg.PathSep("."))
	if err != nil {
		return nil, fmt.Errorf("Fatal error reading config file: %w", err)
	}

	appConfig := defaultConfig
	if err := config.Unpack(&appConfig); err != nil {
		return nil, fmt.Errorf("Unable to parse config file: %v", err)
	}

	loglevel, err := logrus.ParseLevel(appConfig.Loglevel)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse loglevel. %v", err)
	}
	log.SetLevel(loglevel)

	switch appConfig.Logformat {
	case "text": // nothing to do here. Text is the default anyway
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.Fatalf("Unknown log format: %s", appConfig.Logformat)
	}

	if appConfig.DefaultForwarder == nil {
		return nil, fmt.Errorf("defaultForwarder must be specified in configuration file.")
	}
	defaultForwarder, err := NewForwarder(appConfig.DefaultForwarder)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse defaultForwarder: %v", err)
	}

	rules := make([]Rule, len(appConfig.Rules))
	for i, rcfg := range appConfig.Rules {
		rule, err := NewRule(&rcfg)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse rule #%d: %v", i, err)
		}
		rules[i] = *rule
	}

	return &Configuration{
		Bind:             appConfig.Bind,
		Rules:            rules,
		DefaultForwarder: &defaultForwarder,
	}, nil
}
