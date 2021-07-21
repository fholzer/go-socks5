package socks5

import (
	"context"
	"net"
)

// Picker check dest addr and use map[key] of func Dial
type Picker interface {
	Pick(req *Request, ctx context.Context) (context.Context, func(ctx context.Context, network, addr string) (net.Conn, error))
}
