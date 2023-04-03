package main

import (
	"context"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"github.com/partyzanex/powow/pkg/client"
	"github.com/partyzanex/powow/pkg/xnet"
)

var (
	address      = flag.String("address", ":7700", "the server address")
	timeout      = flag.Duration("timeout", time.Second*5, "the request timeout")
	concurrency  = flag.Int("concurrency", 8, "")
	testDuration = flag.Duration("time", time.Minute, "")
)

func main() {
	flag.Parse()

	debugDealer := xnet.NewDialer(&net.Dialer{}, log.Default(), true)
	cli := client.NewClient(*address, client.WithDialer(debugDealer))
	done := make(chan struct{})

	time.AfterFunc(*testDuration, func() {
		close(done)
	})

	wg := sync.WaitGroup{}
	wg.Add(*concurrency)

	for i := 0; i < *concurrency; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				default:
					ctx, cancel := context.WithTimeout(context.Background(), *timeout)

					_, err := cli.GetRandomWisdom(ctx)
					if err != nil {
						log.Println("cannot get random wisdom: ", err)
					}

					cancel()
				}
			}
		}()
	}

	wg.Wait()
}
