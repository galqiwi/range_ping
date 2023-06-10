package connection

import (
	"golang.org/x/net/icmp"
	"net"
)

type PingConnection interface {
	SendEchoRequest(ip net.IP) error
	RecvResponse() (*icmp.Message, net.Addr, error)
	Close() error
}
