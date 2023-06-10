package main

import (
	"fmt"
	"github.com/galqiwi/range_ping/internal/connection"
	"golang.org/x/net/icmp"
	"log"
	"net"
	"time"
)

func timeIt(f func()) {
	begin := time.Now()
	f()
	fmt.Println(time.Since(begin))
}

func Main() error {
	conn, err := connection.NewUDPConnection()
	if err != nil {
		return err
	}

	timeIt(func() {
		err = conn.SendEchoRequest([]byte{1, 1, 1, 1})
	})

	if err != nil {
		return err
	}

	var msg *icmp.Message
	var addr net.Addr

	timeIt(func() {
		msg, addr, err = conn.RecvResponse()
	})

	if err != nil {
		return err
	}

	fmt.Println(msg, addr)

	return nil
}

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}
