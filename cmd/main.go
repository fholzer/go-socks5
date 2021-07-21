package main

import (
	"context"
	"net"
	"log"
	
	"golang.org/x/net/proxy"
	"github.com/fholzer/go-socks5/pkg/socks5"
)

type Picker struct {}

func (p* Picker) Pick(req *socks5.Request, ctx context.Context) (context.Context, func(ctx context.Context, network, addr string) (net.Conn, error)) {
	for _, subnet := range prodSubnets {
		if(subnet.Contains(req.DestAddr.IP)) {
			log.Printf("[INF] Connection to %v will be forwarded via SSH Gateway", req.DestAddr)
			//proxy.SOCKS5(network, addr, nil, nil)).Dial
			return ctx, func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:5656", nil, nil)
				if err != nil {
					return nil, err
				}
				
				return dialer.Dial(network, addr)
			}
		}
	}
	
	log.Printf("[INF] Connection to %v will connect directly", req.DestAddr)
	return ctx, func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial(network, addr)
	}
}

type LogFinalizer struct {
	log *log.Logger
}

func (l *LogFinalizer) Finalize(request *socks5.Request, conn net.Conn, ctx context.Context) error {
	// log.Printf("[INF] Connection from %v to %v ended. rx=%d tx=%d", request.RemoteAddr, request.DestAddr, request.ReqByte, request.RespByte)
	return nil
}


var prodSubnetStrings []string = []string{
}

var prodSubnets []net.IPNet


func main() {
	// create CIDR slots
	prodSubnets = make([]net.IPNet, len(prodSubnetStrings));
	// create CIDRs fron strings
	for i, v := range prodSubnetStrings {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			panic(err)
		}
		prodSubnets[i] = *ipNet
	}
	
	// Create a SOCKS5 server
	conf := &socks5.Config{
		Picker: &Picker{},
		Finalizer: &LogFinalizer{},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "127.0.0.1:5757"); err != nil {
		panic(err)
	}
}