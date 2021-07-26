package main

import (
	"context"
	stdlog "log"
	"net"
	"os"

	"github.com/fholzer/go-socks5/pkg/socks5"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type LogFinalizer struct {
	log *logrus.Logger
}

func (l *LogFinalizer) Finalize(request *socks5.Request, conn net.Conn, ctx context.Context) error {
	log.WithFields(logrus.Fields{
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
	log = &logrus.Logger{
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		Out: ansicolor.NewAnsiColorWriter(os.Stdout),

		// Log as JSON instead of the default ASCII formatter.
		Formatter: &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		},
		Hooks: make(logrus.LevelHooks),

		// Only log the warning severity or above.
		Level: logrus.InfoLevel,
	}

	w := logrus.New().Writer()
	//defer w.Close()
	stdlog.SetOutput(w)
}

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
