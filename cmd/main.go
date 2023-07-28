package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jimmyvo0512/go-libp2p-tutorial/node"
)

func main() {
	// Parse args from command line
	cfg, err := ParseArgs()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Initialize node
	n := node.NewNode()
	if err := n.Start(ctx, cfg.NodePort); err != nil {
		panic(err)
	}
	defer func() {
		if err := n.Shutdown(); err != nil {
			panic(err)
		}
	}()

	// Bootstrap node
	if err := n.Bootstrap(ctx, cfg.BootstrapPeers); err != nil {
		panic(err)
	}

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
