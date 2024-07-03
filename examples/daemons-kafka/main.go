package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Vivino/go-shezmu"
	"github.com/Vivino/go-shezmu/examples/daemons-kafka/daemons"
	"github.com/Vivino/go-shezmu/examples/daemons-kafka/kafka"
	"github.com/Vivino/go-shezmu/server"
	"github.com/Vivino/go-shezmu/stats"
)

var tracker = func(ctx context.Context, name string) (context.Context, func()) {
	return context.Background(), func() {}
}

func main() {
	ctx := context.Background()
	var brokers string

	flag.StringVar(&brokers, "brokers", "127.0.0.1:9092", "Kafka broker addresses separated by space")
	flag.Parse()

	kafka.Initialize(strings.Split(brokers, " "))
	defer kafka.Shutdown()

	statsLogger := stats.NewStdoutLogger(0)
	defer statsLogger.Print()

	statsServer := stats.NewServer()
	server := server.New(6464, statsServer)
	server.Start()

	s := shezmu.Summon()
	s.DaemonStats = stats.NewGroup(statsLogger, statsServer)

	s.AddDaemon(&daemons.NumberPrinter{})
	s.AddDaemon(&daemons.PriceConsumer{})

	s.StartDaemons(ctx, tracker)
	defer s.StopDaemons()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP)
	switch <-sig {
	case syscall.SIGHUP:
		s.StopDaemons()
		s.StartDaemons(ctx, tracker)
	case syscall.SIGINT:
		return
	}
}
