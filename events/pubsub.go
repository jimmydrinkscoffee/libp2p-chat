package events

func NewPubSub() (Publisher, Subscriber) {
	ch := make(chan Event)
	doneCh := make(chan struct{})

	pub := NewPublisher(ch, doneCh)
	go pub.WaitToClose()

	sub := NewSubscriber(ch, doneCh)

	return pub, sub
}
