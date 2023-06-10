package recv_loop

import (
	"fmt"
	"github.com/galqiwi/range_ping/internal/connection"
	"golang.org/x/net/icmp"
	"log"
	"net"
	"sync"
)

type Callback func(*icmp.Message)

type RecvLoop interface {
	AddCallback(ip net.IP, callback Callback)
	CancelCallback(ip net.IP)
}

type recvLoop struct {
	conn connection.PingConnection

	mutex             sync.RWMutex
	callbackByAddress map[[4]byte]Callback
}

func getIPKey(ip net.IP) [4]byte {
	var output [4]byte

	if len(ip) != 4 {
		panic(fmt.Sprintf("invalid ip: %v", ip))
	}

	for i := 0; i < 4; i += 1 {
		output[i] = ip[i]
	}
	return output
}

func (l *recvLoop) getCallback(ip net.IP) Callback {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.callbackByAddress[getIPKey(ip)]
}

func (l *recvLoop) AddCallback(ip net.IP, callback Callback) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.callbackByAddress[getIPKey(ip)] = callback
}

func (l *recvLoop) CancelCallback(ip net.IP) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.callbackByAddress, getIPKey(ip))
}

func (l *recvLoop) run() {
	for {
		msg, addrRaw, err := l.conn.RecvResponse()
		if err != nil {
			log.Fatalf("recv loop error: %v", err.Error())
		}
		addr, ok := addrRaw.(*net.UDPAddr)
		if !ok {
			continue
		}
		callback := l.getCallback(addr.IP)
		if callback == nil {
			continue
		}
		callback(msg)
	}
}

func NewRecvLoop(conn connection.PingConnection) RecvLoop {
	output := recvLoop{
		conn:              conn,
		callbackByAddress: make(map[[4]byte]Callback),
	}

	go output.run()

	return &output
}
