package node

import (
	"context"
	"fmt"
	"time"

	"github.com/jimmyvo0512/go-libp2p-tutorial/util"
	"github.com/libp2p/go-libp2p"
	disc "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const protocolPrefix = "/chat"

type Node interface {
	GetID() peer.ID
	Start(uint16) error
	Bootstrap(context.Context, []multiaddr.Multiaddr) error
	Shutdown() error
}

type node struct {
	host host.Host
	kDht *dht.IpfsDHT
}

var _ Node = (*node)(nil)

func NewNode() Node {
	return &node{
		host: nil,
	}
}

func (n *node) GetID() peer.ID {
	if n.host == nil {
		return util.EmptyString
	}

	return n.host.ID()
}

func (n *node) Start(port uint16) error {
	addr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)

	privKey, _, err := crypto.GenerateEd25519Key(nil)
	if err != nil {
		return err
	}

	host, err := libp2p.New(
		libp2p.ListenAddrStrings(addr),
		libp2p.Identity(privKey),
	)
	if err != nil {
		return err
	}

	n.host = host

	return nil
}

func (n *node) Bootstrap(ctx context.Context, addrs []multiaddr.Multiaddr) error {
	var btNodes []peer.AddrInfo
	for _, addr := range addrs {
		info, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}

		btNodes = append(btNodes, *info)
	}

	kDht, err := dht.New(
		ctx,
		n.host,
		dht.BootstrapPeers(btNodes...),
		dht.ProtocolPrefix(protocolPrefix),
		dht.Mode(dht.ModeAutoServer),
	)
	if err != nil {
		return err
	}

	n.kDht = kDht

	rt := disc.NewRoutingDiscovery(n.kDht)
	disc.Advertise(ctx, rt, protocolPrefix)

	go func() {
		tick := time.NewTicker(time.Second * 5)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				_, err := rt.FindPeers(ctx, protocolPrefix)
				if err != nil {
					continue
				}
			}

		}
	}()

	return nil
}

func (n *node) Shutdown() error {
	if err := n.host.Close(); err != nil {
		return err
	}

	return nil
}
