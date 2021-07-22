package main

import (
	"context"
	"net"
)

type Rule struct {
	subnets   []net.IPNet
	forwarder Forwarder
}

func NewRule(rcfg *ruleConfig) (*Rule, error) {
	// create CIDR array
	subnets := make([]net.IPNet, len(rcfg.Subnets))

	// create CIDRs fron strings
	for i, v := range rcfg.Subnets {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			return nil, err
		}
		subnets[i] = *ipNet
	}

	forwarder, err := NewForwarder(&rcfg.Forwarder)
	if err != nil {
		return nil, err
	}

	rule := &Rule{
		subnets:   subnets,
		forwarder: forwarder,
	}
	return rule, nil
}

func (r *Rule) Match(ip net.IP) bool {
	for _, subnet := range r.subnets {
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}

func (r *Rule) EnrichContext(ctx context.Context) context.Context {
	return r.forwarder.EnrichContext(ctx)
}

func (r *Rule) Forward(ctx context.Context, network, addr string) (net.Conn, error) {
	return r.forwarder.Forward(ctx, network, addr)
}
