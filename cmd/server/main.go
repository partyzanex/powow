package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"

	"github.com/partyzanex/powow/internal/challenge"
	"github.com/partyzanex/powow/internal/quote"
	"github.com/partyzanex/powow/internal/quote/file"
	"github.com/partyzanex/powow/internal/transport"
)

var (
	address               = flag.String("address", ":7700", "the server address")
	challengeMinTime      = flag.Uint("challenge-min-time", 2, "the minimal time factor for PoW challenge")
	challengeMinKeyLength = flag.Uint("challenge-min-key-length", 32, "the minimal key length for PoW challenge")
	challengeMinThreads   = flag.Uint("challenge-min-threads", 1, "the minimal count of threads")
	challengeMinMemory    = flag.Uint("challenge-min-memory", 1, "the minimal memory size")
	challengeTTL          = flag.Duration("challenge-ttl", time.Second*15, "the maximum of challenge response waiting time")
	quotesFilePath        = flag.String("quotes-file-path", "./assets/quotes.txt", "the filepath to database")
	debug                 = flag.Bool("debug", true, "enable/disable debug")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	quotesRepository, err := file.NewRepository(*quotesFilePath)
	if err != nil {
		log.Fatal("cannot create quotes repository:", err)
	}

	service := transport.NewService(
		challenge.NewProvider(&challenge.Config{
			MinTime:      uint32(*challengeMinTime),
			MinKeyLength: uint32(*challengeMinKeyLength),
			MinThreads:   uint8(*challengeMinThreads),
			MinMemory:    uint32(*challengeMinMemory),
			TTL:          *challengeTTL,
		}),
		quote.NewService(quotesRepository),
		*debug,
	)

	n := make(chan os.Signal, 1)
	signal.Notify(n, os.Kill, os.Interrupt)

	go func() {
		<-n

		if closeErr := service.Close(); closeErr != nil {
			log.Println("cannot close service:", closeErr)
		}

		if closeErr := lis.Close(); closeErr != nil {
			log.Println("cannot close listener:", closeErr)
		}
	}()

	log.Println("listen on", *address)

	for {
		conn, err := lis.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			log.Fatal("cannot accept:", err)
		}

		go service.HandleConn(conn)
	}
}
