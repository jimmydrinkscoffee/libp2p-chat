package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jimmyvo0512/go-libp2p-tutorial/node"
)

func main() {
	n := node.NewNode()
	if err := n.Start(0); err != nil {
		panic(err)
	}
	defer func() {
		if err := n.Shutdown(); err != nil {
			panic(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
