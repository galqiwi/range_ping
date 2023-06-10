package main

import (
	"fmt"
	"github.com/galqiwi/range_ping/internal/connection"
	"github.com/galqiwi/range_ping/internal/recv_loop"
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
	defer func() {
		_ = conn.Close()
	}()

	l := recv_loop.NewRecvLoop(conn)
	l.AddCallback(net.IP{1, 1, 1, 1}, func(message *icmp.Message) {
		fmt.Println(message)
	})

	err = conn.SendEchoRequest(net.IP{1, 1, 1, 1})
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	return nil
}

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}
