package socks5

import (
	"context"
	"log"
)

type Finalizer interface {
	Finalize(ctx context.Context) error
}

type LogFinalizer struct {
	log *log.Logger
}

func (l *LogFinalizer) Finalize(ctx context.Context) error {
	l.log.Println(ctx.Value("username"), ctx.Value("raddr"), ctx.Value("daddr"), ctx.Value("request_byte"), ctx.Value("response_byte"))
	return nil
}
