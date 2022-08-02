package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	defaultDuration, _ := time.ParseDuration("10s")
	flag.DurationVar(&timeout, "timeout", defaultDuration, "set timeout")
	flag.Parse()

	host := flag.Arg(0)
	if host == "" {
		log.Fatal("No host given!")
	}

	port := flag.Arg(1)
	if port == "" {
		log.Fatal("No port given!")
	}

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatalf("Got error from client: %v", err)
	}
	defer client.Close()

	sysSignal := make(chan os.Signal, 1)
	defer close(sysSignal)
	signal.Notify(sysSignal, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		_ = client.Send()
	}()

	go func() {
		defer cancel()
		_ = client.Receive()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sysSignal:
			client.Close()
		}
	}
}
