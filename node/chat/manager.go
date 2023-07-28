package chat

import (
	"fmt"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/vbph/libp2p-chat/dtos"
	"github.com/vbph/libp2p-chat/events"
)

type Manager interface {
	JoinRoom(string, string) error
	IsInRoom(string) bool

	SendMessage(string, dtos.Message) error
}

type manager struct {
	nodeID   peer.ID
	ps       *pubsub.PubSub
	kDht     *dht.IpfsDHT
	rooms    map[string]Room
	eventPub events.Publisher
	lock     sync.RWMutex
}

var _ Manager = (*manager)(nil)

func NewManager(nodeID peer.ID, ps *pubsub.PubSub, kDht *dht.IpfsDHT) (Manager, events.Subscriber) {
	eventPub, eventSub := events.NewPubSub()

	m := &manager{
		nodeID:   nodeID,
		ps:       ps,
		kDht:     kDht,
		rooms:    make(map[string]Room),
		eventPub: eventPub,
	}

	go m.advertise()

	return m, eventSub
}

func (m *manager) JoinRoom(r string, name string) error {
	if m.IsInRoom(r) {
		return nil
	}

	return nil
}

func (m *manager) IsInRoom(r string) bool {
	return false
}

func (m *manager) SendMessage(r string, msg dtos.Message) error {
	return nil
}

func (m *manager) advertise() {
	tick := time.NewTicker(time.Second * 5)
	defer tick.Stop()

	for {
		<-tick.C

		func() {
			m.lock.RLock()
			defer m.lock.RUnlock()

			for _, r := range m.rooms {
				m.roomAdvertise(r)
			}
		}()

	}
}

func (m *manager) roomAdvertise(room Room) {
	nn, err := room.GetNickname(m.nodeID)
	if err != nil {
		return
	}

	fmt.Print(nn)
}
