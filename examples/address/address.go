package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ATMackay/go-blocknative/client"
)

func main() {

	// create the base client struct
	cl, err := client.New(context.Background(), client.Opts{
		Scheme: "wss",
		Host:   "api.blocknative.com",
		Path:   "/v0",
		// derive the api key from an environment variable
		// this sets the Client::apiKey field allowing you to retrieve the api key using
		// Client::APIKey
		APIKey: os.Getenv("BLOCKNATIVE_DAPP_ID"),
	})
	if err != nil {
		panic(err)
	}
	// this defers closure of connection and uses proper websockets connection closing semantics
	defer cl.Close()
	// send the initialization message to blocknatives api
	if err := cl.Initialize(client.NewBaseMessageMainnet(cl.APIKey())); err != nil {
		panic(err)
	}
	address := "0xdac17f958d2ee523a2206206994597c13d831ec7"
	// subscribe to events by address
	if err := cl.NewAddressSubscription(address); err != nil {
		panic(err)
	}
	// read messages in a loop
	s := cl.SubscriptionRegistry()[address]
	defer func() { s.Unsubscribe() }()

	for i := 0; i < 10; i++ {
		log.Printf("receive message:\n%v\n", s.PopAll()...)
		time.Sleep(2 * time.Second)
	}
}
