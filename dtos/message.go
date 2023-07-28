package dtos

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Message struct {
	SenderID peer.ID
	SentAt   time.Time
	Content  string
}
