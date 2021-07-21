package main

import (
	"context"
	"net"
	"log"

	"github.com/fholzer/go-socks5/pkg/socks5"
)

type Picker struct {
	rules []Rule
	defaultForwarder Forwarder
}

func (p* Picker) Pick(req *socks5.Request, ctx context.Context) (context.Context, func(ctx context.Context, network, addr string) (net.Conn, error)) {
	for _, rule := range p.rules {
		if(rule.Match(req.DestAddr.IP)) {
			log.Printf("[INF] Connection to %v will be forwarded via SSH Gateway", req.DestAddr)
			//proxy.SOCKS5(network, addr, nil, nil)).Dial
			return ctx, rule.Forward
		}
	}

	log.Printf("[INF] Connection to %v will connect directly", req.DestAddr)
	return ctx, p.defaultForwarder.Forward
}