package main

import (
	"context"
	"net"
	"os"

	"github.com/fholzer/go-socks5/pkg/socks5"
	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
)

type LogFinalizer struct {
	log *log.Logger
}

func (l *LogFinalizer) Finalize(request *socks5.Request, conn net.Conn, ctx context.Context) error {
	log.WithFields(log.Fields{
		"client":         request.RemoteAddr,
		"destination":    request.DestAddr,
		"matchingRuleId": ctx.Value("matchingRuleId"),
		"proxyType":      ctx.Value("proxyType"),
		"proxyAddress":   ctx.Value("proxyAddress"),
		"requestBytes":   request.ReqByte,
		"responseBytes":  request.RespByte,
	}).Debug("Connection closed.")
	return nil
}

const configFileName = "config.yml"

func setupLogging() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func createSocks5Server(appConfig *Configuration) (*socks5.Server, error) {
	// Create a SOCKS5 server
	conf := &socks5.Config{
		Picker: &Picker{
			rules:            appConfig.Rules,
			defaultForwarder: *appConfig.DefaultForwarder,
		},
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
