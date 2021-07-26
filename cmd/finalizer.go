package main

import (
	"context"
	"net"

	"github.com/fholzer/go-socks5/pkg/socks5"
	"github.com/sirupsen/logrus"
)

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
