package main

import (
    "context"
    "fmt"
    "net"

    "golang.org/x/net/proxy"
    "github.com/sirupsen/logrus"
)

type Forwarder interface {
    Forward(ctx context.Context, network, addr string) (net.Conn, error)
}

func NewForwarder(cfg *forwarderConfig) (Forwarder, error) {
    if(cfg.Type == "direct") {
        if(cfg.Address != "") {
            return nil, fmt.Errorf("TEST direct forwarder can't have address!")
        }
        return NewDirectForwarder()
    } else if(cfg.Type == "socks5") {
        return NewSocks5Forwarder(cfg)
    }
    return nil, fmt.Errorf("Unknown forwarder type specified: %s", cfg.Type)
}


type Socks5Forwarder struct {
    dialer proxy.Dialer
    log *logrus.Entry
}

func NewSocks5Forwarder(cfg *forwarderConfig) (*Socks5Forwarder, error) {
    dialer, err := proxy.SOCKS5("tcp", cfg.Address, nil, nil)
    if err != nil {
        return nil, err
    }

    log := logrus.WithFields(logrus.Fields{
        "proxyType": "socks5",
        "proxyAddress": cfg.Address,
    })

    return &Socks5Forwarder{
        dialer: dialer,
        log: log,
    }, nil
}

func (f *Socks5Forwarder) Forward(ctx context.Context, network, addr string) (net.Conn, error) {
    if logrus.IsLevelEnabled(logrus.DebugLevel) {
        f.log.WithFields(logrus.Fields{
            "client": ctx.Value("clientAddr"),
            "destination": addr,
            "matchingRuleId": ctx.Value("matchingRuleId"),
        }).Debug("Forwarding connection via socks5 proxy")
    }

    return f.dialer.Dial(network, addr)
}


type DirectForwarder struct {
    log *logrus.Entry
}

func NewDirectForwarder() (*DirectForwarder, error) {
    log := logrus.WithFields(logrus.Fields{
        "proxyType": "direct",
    })

    return &DirectForwarder{
        log: log,
    }, nil
}

func (f *DirectForwarder) Forward(ctx context.Context, network, addr string) (net.Conn, error) {
    if logrus.IsLevelEnabled(logrus.DebugLevel) {
        f.log.WithFields(logrus.Fields{
            "client": ctx.Value("clientAddr"),
            "destination": addr,
            "matchingRuleId": ctx.Value("matchingRuleId"),
        }).Debug("Forwarding connection directly")
    }

    return net.Dial(network, addr)
}
