package client

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var subscriptionSleep = 5 * time.Second

// Subscription
type Subscription interface {
	Events() chan interface{}
	Unsubscribe()
	Err() chan error
}

func eventLoop(cl *Client, sub Subscription, unsubMsg interface{}) {
	e := sub.Events()
	defer func() {
		sub.Unsubscribe()
		cl.WriteJSON(unsubMsg)
		close(e)
	}()
	for {
		select {
		case <-sub.Err():
			return
		default:
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
			e <- msg
		}
		time.Sleep(subscriptionSleep)
	}
}

type subscription struct {
	key       string // address or txHash
	eventChan chan interface{}
	errChan   chan error
}

func NewSubscription(key string) *subscription {
	return &subscription{key: key, eventChan: make(chan interface{}), errChan: make(chan error, 1)}
}

func (a *subscription) Events() chan interface{} {
	return a.eventChan
}

func (a *subscription) Unsubscribe() {
	a.errChan <- fmt.Errorf("subscription closed")
}

func (a *subscription) Err() chan error {
	return a.errChan
}
