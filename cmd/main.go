package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jimmyvo0512/go-libp2p-tutorial/node"
)

func main() {
	ctx := context.Background()

	n := node.NewNode()
	if err := n.Start(ctx, 0); err != nil {
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
