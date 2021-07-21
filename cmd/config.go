package main

import (
    "fmt"

    ucfg "github.com/elastic/go-ucfg"
    "github.com/elastic/go-ucfg/yaml"
)

type forwarderConfig struct {
	Type string
	Address string
}

type ruleConfig struct {
	Subnets []string
	Forwarder forwarderConfig
}

type rawConfiguration struct {
	Bind string
	Rules []ruleConfig
	DefaultForwarder *forwarderConfig `config:"defaultForwarder"`
}

type Configuration struct {
	Bind string
	Rules []Rule
	DefaultForwarder *Forwarder
}

var (
    defaultConfig = rawConfiguration{
	    Bind: "127.0.0.1:6767",
		Rules: nil,
		DefaultForwarder: nil,
    }
)

func ParseConfig(filename string) (*Configuration, error) {
    config, err := yaml.NewConfigWithFile(filename, ucfg.PathSep("."))
    if  err != nil {
        return nil, fmt.Errorf("Fatal error reading config file: %w\n", err)
    }

    appConfig := defaultConfig
    if err := config.Unpack(&appConfig); err != nil {
        return nil, fmt.Errorf("Unable to parse config file: %v\n", err)
    }

    if appConfig.DefaultForwarder == nil {
        return nil, fmt.Errorf("defaultForwarder must be specified in configuration file.");
    }
	defaultForwarder, err := NewForwarder(appConfig.DefaultForwarder)
    if err != nil {
        return nil, fmt.Errorf("Unable to create defaultForwarder: %v\n", err)
    }

    rules := make([]Rule, len(appConfig.Rules))
    for i, rcfg := range appConfig.Rules {
        rule, err := NewRule(&rcfg)
        if err != nil {
            return nil, fmt.Errorf("Unable to create rule #%d: %v\n", i, err)
        }
        rules[i] = *rule
    }

    return &Configuration{
        Bind: appConfig.Bind,
        Rules: rules,
        DefaultForwarder: &defaultForwarder,
    }, nil
}
