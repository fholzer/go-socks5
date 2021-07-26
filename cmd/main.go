package main

import (
	"net"

	"github.com/fholzer/go-socks5/pkg/socks5"
)

func createSocks5Server(appConfig *Configuration) (*socks5.Server, error) {
	// Create a SOCKS5 server
	conf := &socks5.Config{
		Picker: &Picker{
			rules:            appConfig.Rules,
			defaultForwarder: *appConfig.DefaultForwarder,
		},
		Logger:    log,
		Finalizer: &LogFinalizer{},
	}
	return socks5.New(conf)
}

func main() {
	setupLogging()

	appConfig, err := ParseConfig(configFileName)
	if err != nil {
		log.Fatalf("Error loading configuration file. %v", err)
	}

	server, err := createSocks5Server(appConfig)
	if err != nil {
		log.Panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := ListenAndServe(server, "tcp", appConfig.Bind); err != nil {
		log.Fatal(err)
	}
}

func ListenAndServe(s *socks5.Server, network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalf("Error binding to %s. %v", addr, err)
	}
	log.Info("Server running and waiting for connections...")
	return s.Serve(l)
}
