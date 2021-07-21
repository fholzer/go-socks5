package main

import (
	"context"
	"net"
	"log"
	
	"github.com/fholzer/go-socks5/pkg/socks5"
)

type LogFinalizer struct {
	log *log.Logger
}

func (l *LogFinalizer) Finalize(request *socks5.Request, conn net.Conn, ctx context.Context) error {
	// log.Printf("[INF] Connection from %v to %v ended. rx=%d tx=%d", request.RemoteAddr, request.DestAddr, request.ReqByte, request.RespByte)
	return nil
}

const configFileName = "config.yml"

func main() {
	appConfig, err := ParseConfig(configFileName)
	if err != nil {
		panic(err)
	}
	
	// Create a SOCKS5 server
	conf := &socks5.Config{
		Picker: &Picker{
			rules: appConfig.Rules,
			defaultForwarder: *appConfig.DefaultForwarder,
		},
		Finalizer: &LogFinalizer{},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", appConfig.Bind); err != nil {
		panic(err)
	}
}