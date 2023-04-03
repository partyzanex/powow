package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/partyzanex/powow/pkg/client"
	"github.com/partyzanex/powow/pkg/xnet"
)

var (
	address = flag.String("address", ":7700", "the server address")
	debug   = flag.Bool("debug", true, "enable/disable debug")
	timeout = flag.Duration("timeout", time.Second*5, "the request timeout")
)

func main() {
	flag.Parse()

	var options []client.Option

	if *debug {
		debugDealer := xnet.NewDialer(&net.Dialer{}, log.Default(), *debug)
		options = append(options, client.WithDialer(debugDealer))
	}

	cli := client.NewClient(*address, options...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := make(chan os.Signal, 1)
	signal.Notify(n, os.Kill, os.Interrupt)

	go func() {
		<-n
		cancel()
	}()

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, *timeout)
	defer timeoutCancel()

	quote, err := cli.GetRandomWisdom(timeoutCtx)
	if err != nil {
		log.Fatal("cannot get random wisdom: ", err)
	}

	fmt.Printf("Quote: %s\nAuthor: %s\n", quote.Content, quote.Author)
}
