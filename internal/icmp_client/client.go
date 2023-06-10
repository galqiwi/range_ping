package icmp_client

import (
	"errors"
	"github.com/galqiwi/range_ping/internal/connection"
	"github.com/galqiwi/range_ping/internal/recv_loop"
	"golang.org/x/net/icmp"
	"net"
	"time"
)

var TimeoutErr = errors.New("timeouted")

type ICMPClient interface {
	SendRequest(ip net.IP) (*icmp.Message, error)
}

type icmpClient struct {
	timeout time.Duration
	conn    connection.PingConnection
	rl      recv_loop.RecvLoop
}

func (c *icmpClient) SendRequest(ip net.IP) (*icmp.Message, error) {
	respChannel := make(chan *icmp.Message, 1)
	c.rl.AddCallback(ip, func(message *icmp.Message) {
		respChannel <- message
	})
	defer c.rl.CancelCallback(ip)
	err := c.conn.SendEchoRequest(ip)
	if err != nil {
		return nil, err
	}

	select {
	case resp := <-respChannel:
		return resp, nil
	case <-time.After(c.timeout):
		return nil, TimeoutErr
	}
}

func NewICMPClient(timeout time.Duration) (ICMPClient, error) {
	conn, err := connection.NewUDPConnection()
	if err != nil {
		return nil, err
	}
	rl := recv_loop.NewRecvLoop(conn)
	return &icmpClient{
		conn:    conn,
		rl:      rl,
		timeout: timeout,
	}, nil
}
