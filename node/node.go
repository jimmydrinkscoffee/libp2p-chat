package node

import (
	"fmt"

	"github.com/jimmyvo0512/go-libp2p-tutorial/util"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Node interface {
	GetID() peer.ID
	Start(port uint16) error
	Shutdown() error
}

type node struct {
	host host.Host
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

func (n *node) Shutdown() error {
	if err := n.host.Close(); err != nil {
		return err
	}

	return nil
}
