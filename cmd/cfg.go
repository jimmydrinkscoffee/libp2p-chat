package main

import (
	"flag"
	"strings"

	"github.com/jimmyvo0512/go-libp2p-tutorial/util"
	"github.com/multiformats/go-multiaddr"
)

type Cfg struct {
	NodePort       uint16
	BootstrapPeers []multiaddr.Multiaddr
}

func ParseArgs() (Cfg, error) {
	nodePort := flag.Uint("nodePort", 0, "Node port")
	bootstrapPeers := flag.String("bootstrapPeers", util.EmptyString, "Bootstrap peer addresses separated by comma")

	flag.Parse()

	if *nodePort == 0 {
		return Cfg{}, flag.ErrHelp
	}

	var btPeers []multiaddr.Multiaddr
	if *bootstrapPeers != util.EmptyString {
		for _, p := range strings.Split(*bootstrapPeers, ",") {
			addr, err := multiaddr.NewMultiaddr(p)
			if err != nil {
				return Cfg{}, err
			}

			btPeers = append(btPeers, addr)
		}
	}

	return Cfg{
		NodePort:       uint16(*nodePort),
		BootstrapPeers: btPeers,
	}, nil
}
