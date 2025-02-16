package client

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// Opts provides configuration over the websocket connection
type Opts struct {
	Scheme               string
	Host                 string
	Path                 string
	APIKey               string
	PrintConnectResponse bool
}

// ConnectResponse is the message we receive when opening a connection to the API
type ConnectResponse struct {
	ConnectionID  string `json:"connectionId"`
	ServerVersion string `json:"serverVersion"`
	ShowUX        bool   `json:"showUX"`
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	Version       int    `json:"version"`
}

// Client wraps gorilla websocket connections
type Client struct {
	conn                 *websocket.Conn
	ctx                  context.Context
	cancel               context.CancelFunc
	initMsg              BaseMessage // used to resend the initialization msg if connection drops
	apiKey               string
	mtx                  sync.RWMutex
	subscriptionRegistry map[string]Subscription
}

// New returns a new blocknative websocket client
func New(ctx context.Context, opts Opts) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	u := url.URL{
		Scheme: opts.Scheme,
		Host:   opts.Host,
		Path:   opts.Path,
	}
	c, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		cancel()
		return nil, err
	}
	// this checks out connection to blocknative's api and makes sure that we connected properly
	var out ConnectResponse
	if err := c.ReadJSON(&out); err != nil {
		cancel()
		return nil, err
	}
	if out.Status != "ok" {
		cancel()
		return nil, fmt.Errorf("failed to initialize websockets connection reason: %v", out.Reason)
	}
	if opts.PrintConnectResponse {
		log.Printf("%+v\n", out)
	}
	return &Client{conn: c, ctx: ctx, cancel: cancel, apiKey: opts.APIKey, subscriptionRegistry: make(map[string]Subscription)}, nil
}

// Initialize is used to handle blocknative websockets api initialization
// note we set CategoryCode and EventCode ourselves.
func (c *Client) Initialize(msg BaseMessage) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	msg.Version = "1"
	msg.CategoryCode = "initialize"
	msg.EventCode = "checkDappId"
	c.initMsg = msg
	if err := c.conn.WriteJSON(&msg); err != nil {
		return err
	}
	var out ConnectResponse
	err := c.conn.ReadJSON(&out)
	if err != nil {
		return err
	}
	if out.Status != "ok" {
		return fmt.Errorf("failed to initialize api connection reason:%v", out.Reason)
	}
	return nil
}

// APIKey returns the api key being used by the client
func (c *Client) APIKey() string {
	return c.apiKey
}

// SubscriptionRegistry returns the chached subscription map
func (c *Client) SubscriptionRegistry() map[string]Subscription {
	return c.subscriptionRegistry
}

// NewEventSubscription creates an event subscription.
func (c *Client) NewEventSubscription(msg Configuration) error {
	if err := c.WriteJSON(&msg); err != nil {
		return err
	}
	var out ConnectResponse
	if err := c.ReadJSON(&out); err != nil {
		return err
	}
	if out.Status != "ok" {
		return fmt.Errorf("failed to create subscription reason:%v", out.Reason)
	}

	key := msg.Scope
	c.subscriptionRegistry[key] = NewSubscription(key)
	go eventLoop(c, c.subscriptionRegistry[key], NewEventUnsubscribe(c.initMsg, msg.Config))

	return nil
}

// NewAddressSubscription creates a new subscription to blocknative
// for monitoring address activity. the subscription is added to the
// client's subscription registry which contains an event channel
// for watched events provided by blocknative servers
func (c *Client) NewAddressSubscription(address string) error {
	if err := c.WriteJSON(NewAddressSubscribe(
		c.initMsg,
		address,
	)); err != nil {
		return err
	}
	c.subscriptionRegistry[address] = NewSubscription(address)
	go eventLoop(c, c.subscriptionRegistry[address], NewAddressUnsubscribe(c.initMsg, address))
	return nil
}

// NewTransactionSubscription creates a new subscription for monitoring
// a transaction by supplied transaction ID
func (c *Client) NewTransactionSubscription(txHash string) error {
	if err := c.WriteJSON(NewTxSubscribe(
		c.initMsg,
		txHash,
	)); err != nil {
		return err
	}
	c.subscriptionRegistry[txHash] = NewSubscription(txHash)
	go eventLoop(c, c.subscriptionRegistry[txHash], NewTxUnsubscribe(c.initMsg, txHash))
	return nil
}

// KillSubscription unsubscribes from the subscription in the registry
// at the index supplied. Killing the subscription releases the resource for the Client
// as well notifying block native servers to not monitor for events tracked by the subscription
func (c *Client) KillSubscription(key string) {
	sub, ok := c.subscriptionRegistry[key]
	if !ok {
		// no subscription found
		return
	}
	sub.Unsubscribe()
	delete(c.subscriptionRegistry, key)
}

// ReadJSON is a wrapper around Conn:ReadJSON
func (c *Client) ReadJSON(out interface{}) error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.conn.ReadJSON(out)
}

// WriteJSON is a wrapper around Conn:WriteJSON
func (c *Client) WriteJSON(out interface{}) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.conn.WriteJSON(out)
}

// Close is used to terminate our websocket client
func (c *Client) Close() error {
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	c.cancel()
	return err
}
