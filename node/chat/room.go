package chat

import (
	"errors"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/vbph/libp2p-chat/util"
)

type Room interface {
	GetNickname(peer.ID) (string, error)
}

type participant struct {
	ID       peer.ID
	Nickname string
}

type room struct {
	ptcps map[peer.ID]*participant
	lock  sync.RWMutex
}

var _ Room = (*room)(nil)

func (r *room) GetNickname(nodeID peer.ID) (string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if p, ok := r.ptcps[nodeID]; ok {
		return p.Nickname, nil
	}

	return util.EmptyString, errors.New(util.EmptyString)
}
