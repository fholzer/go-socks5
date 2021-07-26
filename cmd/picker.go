package main

import (
	"context"
	"net"

	"github.com/fholzer/go-socks5/pkg/socks5"
	"github.com/sirupsen/logrus"
)

type Picker struct {
	rules            []Rule
	defaultForwarder Forwarder
}

func (p *Picker) Pick(req *socks5.Request, ctx context.Context) (context.Context, func(ctx context.Context, network, addr string) (net.Conn, error)) {
	var logentry *logrus.Entry
	if log.IsLevelEnabled(logrus.TraceLevel) {
		logentry = log.WithFields(logrus.Fields{
			"client":      req.RemoteAddr,
			"destination": req.DestAddr,
		})
		logentry.Trace("Starting rule processing.")
	}

	ctx = context.WithValue(ctx, "clientAddr", req.RemoteAddr)

	for i, rule := range p.rules {
		if rule.Match(req.DestAddr.IP) {
			if logentry != nil {
				logentry.WithField("matchingRuleId", i).Tracef("Rule %d matches.", i)
			}
			//proxy.SOCKS5(network, addr, nil, nil)).Dial
			ctx = context.WithValue(ctx, "matchingRuleId", i)
			ctx = rule.EnrichContext(ctx)
			return ctx, rule.Forward
		}
		if logentry != nil {
			logentry.WithField("matchingRuleId", i).Tracef("Rule %d doesn't matche.", i)
		}
	}

	if logentry != nil {
		logentry.Tracef("Using fallback forwarder.")
	}
	ctx = context.WithValue(ctx, "matchingRuleId", -1)
	ctx = p.defaultForwarder.EnrichContext(ctx)
	return ctx, p.defaultForwarder.Forward
}
