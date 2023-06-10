package main

import (
	"fmt"
	"github.com/galqiwi/range_ping/internal/icmp_client"
	"log"
	"net"
	"sync"
	"time"
)

func timeIt(f func()) {
	begin := time.Now()
	f()
	fmt.Println(time.Since(begin))
}

func Main() error {
	client, err := icmp_client.NewICMPClient(time.Second * 3)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < 256; i += 1 {
		i := i
		wg.Add(1)
		go func() {
			resp, err := client.SendRequest(net.IP{1, 1, 1, byte(i)})
			fmt.Println(i, resp, err)
			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}
