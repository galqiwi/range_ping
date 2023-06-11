package main

import (
	"fmt"
	"github.com/galqiwi/range_ping/internal/icmp_client"
	"github.com/galqiwi/range_ping/internal/mask_iterator"
	"github.com/paulbellamy/ratecounter"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var requestCounter = ratecounter.NewRateCounter(1 * time.Second)
var successCounter = ratecounter.NewRateCounter(1 * time.Second)
var retryCounter = ratecounter.NewRateCounter(1 * time.Second)

func worker(client icmp_client.ICMPClient, ips chan net.IP) {
	for ip := range ips {
		var err error
		for {
			retryCounter.Incr(1)
			_, err = client.SendRequest(ip)
			if err != nil && strings.Contains(err.Error(), "no buffer space available") {
				time.Sleep(time.Millisecond * time.Duration(rand.Int()%2000))
				continue
			}
			break
		}
		requestCounter.Incr(1)
		if err == nil {
			fmt.Println(ip)
			successCounter.Incr(1)
		}
	}
}

func Main() error {
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(retryCounter.Rate(), requestCounter.Rate(), successCounter.Rate())
		}
	}()

	ips, err := mask_iterator.IPGenerator("1.1.1.1/16")
	if err != nil {
		return err
	}

	client, err := icmp_client.NewICMPClient(time.Second * 3)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < 10000; i += 1 {
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
