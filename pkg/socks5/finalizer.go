package socks5

import (
	"context"
	"net"

	"github.com/fholzer/go-socks5/pkg/axe"
)

type Finalizer interface {
	Finalize(request *Request, conn net.Conn, ctx context.Context) error
}

type LogFinalizer struct {
	log axe.Logger
}

func (l *LogFinalizer) Finalize(request *Request, conn net.Conn, ctx context.Context) error {
	return nil
}
