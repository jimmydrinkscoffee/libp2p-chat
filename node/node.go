package node

import (
	"context"
	"fmt"
	"time"

	"github.com/jimmyvo0512/go-libp2p-tutorial/node/chat"
	"github.com/jimmyvo0512/go-libp2p-tutorial/util"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	discRt "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discUtil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

const protocolPrefix = "/chat"

type Node interface {
	GetID() peer.ID
	Start(context.Context, uint16) error
	Bootstrap(context.Context, []multiaddr.Multiaddr) error
	Shutdown() error
}

type node struct {
	host    host.Host
	kDht    *dht.IpfsDHT
	ps      *pubsub.PubSub
	chatMgr chat.Manager
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

func (n *node) Start(ctx context.Context, port uint16) error {
	if err := n.initHost(port); err != nil {
		return err
	}

	if err := n.initPubSub(ctx); err != nil {
		return err
	}

	return nil
}

func (n *node) Bootstrap(ctx context.Context, addrs []multiaddr.Multiaddr) error {
	var btPeers []peer.AddrInfo
	for _, addr := range addrs {
		info, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}

		btPeers = append(btPeers, *info)
	}

	kDht, err := dht.New(
		ctx,
		n.host,
		dht.BootstrapPeers(btPeers...),
		dht.ProtocolPrefix(protocolPrefix),
		dht.Mode(dht.ModeAutoServer),
	)
	if err != nil {
		return err
	}

	n.kDht = kDht

	rt := discRt.NewRoutingDiscovery(n.kDht)
	discUtil.Advertise(ctx, rt, protocolPrefix)

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

func (n *node) initHost(port uint16) error {
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

func (n *node) initPubSub(ctx context.Context) error {
	ps, err := pubsub.NewGossipSub(
		ctx,
		n.host,
		pubsub.WithMessageSignaturePolicy(pubsub.StrictSign),
	)
	if err != nil {
		return err
	}

	n.ps = ps
	return nil
}
