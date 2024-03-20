package channel

import (
	"errors"
	"sync"
)

type Broker struct {
	mutex sync.RWMutex
	chans []chan Msg
}

func (b *Broker) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for i := 0; i < len(b.chans); i++ {
		select {
		case b.chans[i] <- m:
		default:
			return errors.New("消息队列已满")
		}
	}
	return nil
}

func (b *Broker) Subscribe(capacity int) (<-chan Msg, error) {
	res := make(chan Msg, capacity)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.chans = append(b.chans, res)
	return res, nil
}

func (b *Broker) Close() error {
	b.mutex.Lock()
	chans := b.chans
	b.chans = nil
	defer b.mutex.Unlock()
	// 避免重复close问题
	for _, ch := range chans {
		close(ch)
	}
	return nil
}

type Msg struct {
	Content string
}

type BrokerV2 struct {
	mutex    sync.RWMutex
	consumer []func(msg Msg)
}

func (b *BrokerV2) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, c := range b.consumer {
		c(m)
	}
	return nil
}

func (b *BrokerV2) Subscribe(cb func(s Msg)) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.consumer = append(b.consumer, cb)
	return nil
}
