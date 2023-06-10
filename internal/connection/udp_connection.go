package connection

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"sync"
)

const bufferSize = 65535 * 2

type udpConnection struct {
	internal   net.PacketConn
	recvBuffer []byte

	// just in case, I don't feel confident in concurrent
	// access to net.PacketConn
	sendMutex sync.Mutex
	recvMutex sync.Mutex
}

func (c *udpConnection) Close() error {
	return c.internal.Close()
}

func (c *udpConnection) SendEchoRequest(ip net.IP) error {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte{},
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return err
	}

	if _, err := c.internal.WriteTo(wb, &net.UDPAddr{IP: ip, Port: 80}); err != nil {
		return err
	}

	return nil
}

func (c *udpConnection) RecvResponse() (*icmp.Message, net.Addr, error) {
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()

	n, peer, err := c.internal.ReadFrom(c.recvBuffer)
	if err != nil {
		return nil, nil, err
	}
	if n == bufferSize {
		panic("internal error: buffer is not large enough")
	}

	output, err := icmp.ParseMessage(1, c.recvBuffer[:n])
	if err != nil {
		return nil, nil, err
	}

	return output, peer, nil
}

func NewUDPConnection() (PingConnection, error) {
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return nil, err
	}

	return &udpConnection{
		internal:   conn,
		recvBuffer: make([]byte, bufferSize),
	}, nil
}
