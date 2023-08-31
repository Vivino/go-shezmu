package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Vivino/go-shezmu"
	"github.com/Vivino/go-shezmu/example/daemons"
	"github.com/Vivino/go-shezmu/example/kafka"
	"github.com/Vivino/go-shezmu/server"
	"github.com/Vivino/go-shezmu/stats"
)

func main() {
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
	s.Subscriber = kafka.Subscriber{}
	s.DaemonStats = stats.NewGroup(statsLogger, statsServer)

	s.AddDaemon(&daemons.NumberPrinter{})
	s.AddDaemon(&daemons.PriceConsumer{})

	s.StartDaemons()
	defer s.StopDaemons()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP)
	switch <-sig {
	case syscall.SIGHUP:
		s.StopDaemons()
		s.StartDaemons()
	case syscall.SIGINT:
		return
	}
}
