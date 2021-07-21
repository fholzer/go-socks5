package main

import (
    "context"
    "fmt"
    "net"

    "golang.org/x/net/proxy"
)

type Forwarder interface {
    Forward(ctx context.Context, network, addr string) (net.Conn, error)
}

func NewForwarder(cfg *forwarderConfig) (Forwarder, error) {
    if(cfg.Type == "direct") {
        if(cfg.Address != "") {
            panic("TEST direct forwarder can't have address!")
        }
        return NewDirectForwarder()
    } else if(cfg.Type == "socks5") {
        return NewSocks5Forwarder(cfg)
    }
    panic(fmt.Sprintf("Unknown forwarder type specified: %s", cfg.Type))
}


type Socks5Forwarder struct {
    dialer proxy.Dialer
}

func NewSocks5Forwarder(cfg *forwarderConfig) (*Socks5Forwarder, error) {
    dialer, err := proxy.SOCKS5("tcp", cfg.Address, nil, nil)
    if err != nil {
        return nil, err
    }

    return &Socks5Forwarder{
        dialer: dialer,
    }, nil
}

func (f *Socks5Forwarder) Forward(ctx context.Context, network, addr string) (net.Conn, error) {
    return f.dialer.Dial(network, addr)
}


type DirectForwarder struct {}

func NewDirectForwarder() (*DirectForwarder, error) {
    return &DirectForwarder{}, nil
}

func (f *DirectForwarder) Forward(ctx context.Context, network, addr string) (net.Conn, error) {
    return net.Dial(network, addr)
}
