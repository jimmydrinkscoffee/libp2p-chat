package events

import (
	"errors"

	"github.com/vbph/libp2p-chat/util"
)

type Subscriber interface {
	Next() (Event, error)
	Close()
}

type subscriber struct {
	revcCh   <-chan Event
	doneCh   chan<- struct{}
	isClosed bool
}

var _ Subscriber = (*subscriber)(nil)

func NewSubscriber(revcCh <-chan Event, doneCh chan<- struct{}) Subscriber {
	return &subscriber{
		revcCh:   revcCh,
		doneCh:   doneCh,
		isClosed: false,
	}
}

func (s *subscriber) Next() (Event, error) {
	if s.isClosed {
		return nil, errors.New(util.EmptyString)
	}

	return <-s.revcCh, nil
}

func (s *subscriber) Close() {
	s.doneCh <- struct{}{}
	close(s.doneCh)
	s.isClosed = true
}
