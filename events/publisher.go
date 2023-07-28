package events

import (
	"errors"
	"sync"

	"github.com/vbph/libp2p-chat/util"
)

type Publisher interface {
	Publish(Event) error
	WaitToClose()
	IsClosed() bool
}

type publisher struct {
	sendCh   chan<- Event
	doneCh   <-chan struct{}
	isClosed bool
	lock     sync.RWMutex
}

var _ Publisher = (*publisher)(nil)

func NewPublisher(sendCh chan<- Event, doneCh <-chan struct{}) Publisher {
	return &publisher{
		sendCh:   sendCh,
		doneCh:   doneCh,
		isClosed: false,
	}
}

func (p *publisher) Publish(event Event) error {
	if p.IsClosed() {
		return errors.New(util.EmptyString)
	}

	select {
	case <-p.doneCh:
	case p.sendCh <- event:
		return errors.New(util.EmptyString)
	}

	return nil
}

func (p *publisher) WaitToClose() {
	<-p.doneCh

	p.lock.Lock()
	defer p.lock.Unlock()

	p.isClosed = true
	close(p.sendCh)
}

func (p *publisher) IsClosed() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.isClosed
}
