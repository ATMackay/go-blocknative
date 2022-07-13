package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Subscription implements
type Subscription interface {
	Push(interface{})
	Pop() interface{}
	PopAll() []interface{}
	Unsubscribe()
	Err() chan error
}

func eventLoop(cl *Client, sub Subscription, unsubMsg interface{}) {
	defer func() {
		sub.Unsubscribe()
		cl.WriteJSON(unsubMsg)
	}()
	for {
		select {
		case <-sub.Err():
			return
		default:
			time.Sleep(5 * time.Second)
			var msg EthTxPayload
			if err := cl.ReadJSON(&msg); err != nil {
				if err := cl.ReadJSON(msg); err != nil {
					if e, ok := err.(*websocket.CloseError); ok {
						if e.Code != 1000 {
							sub.Err() <- fmt.Errorf("websocket close error: %v", err)
							return
						}
					}
					break
				} else {
					break
				}
			}
			sub.Push(msg)
		}
	}
}

type EventSubscription struct {
	key     string // address or txHash
	history *MsgHistory
	errChan chan error
}

func NewEventSubscription(key string) EventSubscription {
	return EventSubscription{key: key, history: &MsgHistory{mx: sync.RWMutex{}, buffer: make([]interface{}, 0)}, errChan: make(chan error, 1)}
}

func (a EventSubscription) Push(msg interface{}) {
	a.history.Push(msg)
}

func (a EventSubscription) Pop() interface{} {
	return a.history.Pop()
}

func (a EventSubscription) PopAll() []interface{} {
	return a.history.PopAll()
}

func (a EventSubscription) Unsubscribe() {
	a.errChan <- fmt.Errorf("subscription closed")
}

func (a EventSubscription) Err() chan error {
	return a.errChan
}

/*
// NewSubscription runs a producer function as a subscription in a new goroutine. The
// channel given to the producer is closed when Unsubscribe is called. If fn returns an
// error, it is sent on the subscription's error channel.
func NewSubscription(producer func(<-chan struct{}) error) Subscription {
	s := &funcSub{unsub: make(chan struct{}), err: make(chan error, 1)}
	go func() {
		defer close(s.err)
		err := producer(s.unsub)
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.unsubscribed {
			if err != nil {
				s.err <- err
			}
			s.unsubscribed = true
		}
	}()
	return s
}

type funcSub struct {
	unsub        chan struct{}
	err          chan error
	mu           sync.Mutex
	unsubscribed bool
}

func (s *funcSub) Unsubscribe() {
	s.mu.Lock()
	if s.unsubscribed {
		s.mu.Unlock()
		return
	}
	s.unsubscribed = true
	close(s.unsub)
	s.mu.Unlock()
	// Wait for producer shutdown.
	<-s.err
}

func (s *funcSub) Err() <-chan error {
	return s.err
}

*/
