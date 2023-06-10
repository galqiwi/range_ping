package main

import (
	"fmt"
	"github.com/galqiwi/range_ping/internal/icmp_client"
	"github.com/galqiwi/range_ping/internal/mask_iterator"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func timeIt(f func()) {
	begin := time.Now()
	f()
	fmt.Println(time.Since(begin))
}

func worker(client icmp_client.ICMPClient, ips chan net.IP) {
	for ip := range ips {
		for {
			msg, err := client.SendRequest(ip)
			if err != nil && strings.Contains(err.Error(), "no buffer space available") {
				time.Sleep(time.Second)
				continue
			}
			fmt.Println(ip, msg, err)
		}
	}
}

func Main() error {
	ips, err := mask_iterator.IPGenerator("1.1.1.1/8")
	if err != nil {
		return err
	}

	client, err := icmp_client.NewICMPClient(time.Second * 3)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < 5000; i += 1 {
		wg.Add(1)
		go func() {
			worker(client, ips)
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
