package socks5

import (
	"context"
	"log"
	"net"
)

type Finalizer interface {
	Finalize(request *Request, conn net.Conn, ctx context.Context) error
}

type LogFinalizer struct {
	log *log.Logger
}

func (l *LogFinalizer) Finalize(request *Request, conn net.Conn, ctx context.Context) error {
	return nil
}
