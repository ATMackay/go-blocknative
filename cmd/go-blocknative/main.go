package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ATMackay/go-blocknative/client"
	"github.com/urfave/cli/v2"
)

var (
	apiClient *client.Client
)

func main() {
	app := cli.NewApp()
	app.Name = "go-blocknative"
	app.Usage = "cli for interacting with blocknative api"
	app.Before = func(c *cli.Context) (err error) {
		apiClient, err = client.New(c.Context, client.Opts{
			Scheme: c.String("scheme"),
			Host:   c.String("host"),
			Path:   c.String("api.path"),
			APIKey: c.String("api.key"),
		})
		if err != nil {
			return
		}
		err = apiClient.Initialize(client.NewBaseMessageMainnet(c.String("api.key")))
		return
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "api.key",
			EnvVars: []string{"BLOCKNATIVE_DAPP_ID"},
			Usage:   "blocknative api key",
		},
		&cli.StringFlag{
			Name:  "address",
			Usage: "address to use when subscribing to events",
			Value: "0xfa6de2697D59E88Ed7Fc4dFE5A33daC43565ea41",
		},
		&cli.StringFlag{
			Name:  "tx.hash",
			Usage: "transaction hash to use when subscribing to events",
		},
		&cli.StringFlag{
			Name:  "scheme",
			Usage: "connection scheme to use",
			Value: "wss",
		},
		&cli.StringFlag{
			Name:  "host",
			Usage: "host to connect to",
			Value: "api.blocknative.com",
		},
		&cli.StringFlag{
			Name:  "api.path",
			Usage: "api path to use",
			Value: "/v0",
		},
	}
	app.Commands = cli.Commands{
		&cli.Command{
			Name:    "subscribe",
			Aliases: []string{"sub"},
			Usage:   "event subscription commands",
			Subcommands: cli.Commands{
				&cli.Command{
					Name:  "address",
					Usage: "subscribe to events based on address",
					Action: func(c *cli.Context) error {
						address := c.String("address")
						apiClient.NewAddressSubscription(address)
						s := apiClient.SubscriptionRegistry()[address]
						eventChan := s.Events()
						go func() {
							for {
								e, ok := <-eventChan
								if !ok {
									break
								}
								ev := e.(client.EthTxPayload)
								jev, _ := json.Marshal(ev)
								log.Printf("receive message:\n%v\n", string(jev))
								time.Sleep(5 * time.Second)
							}
						}()
						// start the signal handler
						signalChan := make(chan os.Signal, 1)
						signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
						sig := <-signalChan
						log.Printf("received shutdown signal '%v', unsubscribing...\n", sig)
						apiClient.KillSubscription(address)
						log.Printf("bye!\n")
						return nil
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
