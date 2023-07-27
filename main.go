package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	multiaddr "github.com/multiformats/go-multiaddr"
)

const pid protocol.ID = "/chat/1.0.0"

func handleMessage(s network.Stream) {
	defer s.Close()

	buf := make([]byte, 1024)
	n, err := s.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received message: %s\n", string(buf[:n]))
}

func main() {
	// init node
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	// init ping service
	// pingService := &ping.PingService{Host: node}
	// node.SetStreamHandler(ping.ID, pingService.PingHandler)

	node.SetStreamHandler(pid, handleMessage)

	// print node address
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("libp2p node address:", addrs[0])

	if len(os.Args) > 1 {
		addr, err := multiaddr.NewMultiaddr(os.Args[1])
		if err != nil {
			panic(err)
		}

		peer, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			panic(err)
		}

		if err := node.Connect(context.Background(), *peer); err != nil {
			panic(err)
		}

		s, err := node.NewStream(context.Background(), peer.ID, pid)
		if err != nil {
			panic(err)
		}
		defer s.Close()

		msg := "Hello, world!"
		_, err = s.Write([]byte(msg))
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("Sent message: %s\n", msg)
		}

		// fmt.Println("sending 5 ping messages to", addr)

		// ch := pingService.Ping(context.Background(), peer.ID)
		// for i := 0; i < 5; i++ {
		// 	res := <-ch
		// 	fmt.Println("pinged", addr, "in", res.RTT)
		// }
	} else {
		// wait for interrupt signal
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Received signal, shutting down...")
	}

	// shutdown node
	if err := node.Close(); err != nil {
		panic(err)
	}
}
