package node

import (
	"context"

	"github.com/jimmyvo0512/go-libp2p-tutorial/util"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Node interface {
	GetID() peer.ID
	Start(ctx context.Context, port uint16) error
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

func (n *node) Start(ctx context.Context, port uint16) error {
	return nil
}

func (n *node) Shutdown() error {
	if err := n.host.Close(); err != nil {
		return err
	}

	return nil
}
