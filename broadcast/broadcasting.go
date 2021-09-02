package broadcast

import (
	"context"
	"sync"
)

type InitFunc func(context.Context, *sync.WaitGroup) chan<- Message

type FilterFunc func(Message) bool

type Message interface{}

type Broadcaster struct {
	listeners []chan<- Message
	filters   []FilterFunc
	waitgroup *sync.WaitGroup
	mutex     sync.Mutex
}

func NewBroadcaster() *Broadcaster {
	broadcaster := new(Broadcaster)
	broadcaster.listeners = make([]chan<- Message, 0)
	broadcaster.filters = make([]FilterFunc, 0)
	broadcaster.waitgroup = new(sync.WaitGroup)
	return broadcaster
}

func (b *Broadcaster) Listeners(ctx context.Context, listeners ...InitFunc) *Broadcaster {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.waitgroup.Add(len(listeners))
	for _, init := range listeners {
		b.listeners = append(b.listeners, init(ctx, b.waitgroup))
	}
	return b
}

func (b *Broadcaster) Filters(filters ...FilterFunc) *Broadcaster {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.filters = append(b.filters, filters...)
	return b
}

func (b *Broadcaster) Broadcast(msg Message) *Broadcaster {
	for _, filter := range b.filters {
		if filter(msg) {
			return b
		}
	}
	for _, ch := range b.listeners {
		ch <- msg
	}
	return b
}

func (b *Broadcaster) Wait() {
	b.waitgroup.Wait()
}
