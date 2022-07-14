package client

import (
	"fmt"
	"os"
	"time"
)

// NetName converts chain ID to network name (string)
func NetName(id int64) (string, error) {
	var netName string
	switch id {
	case 1:
		netName = "main"
	case 3:
		netName = "ropsten"
	case 4:
		netName = "rinkeby"
	case 5:
		netName = "goerli"
	case 42:
		netName = "kovan"
	case 56:
		netName = "bsc-main"
	case 100:
		netName = "xdai"
	case 137:
		netName = "matic-main"
	case 250:
		netName = "fantom-main"
	default:
		return "", fmt.Errorf("network not supported id: %v", id)
	}
	return netName, nil
}

// NewBaseMessageMainnet returns a base message suitable for mainnet usage
func NewBaseMessageMainnet(apiKey string) BaseMessage {
	if apiKey == "" {
		apiKey = os.Getenv("BLOCKNATIVE_DAPP_ID")
	}
	return BaseMessage{
		Timestamp: time.Now(),
		DappID:    apiKey,
		Blockchain: Blockchain{
			System:  "ethereum",
			Network: "main",
		},
	}
}

// NewBaseMessage returns a base message for the supplied network ID
func NewBaseMessage(apiKey string, netID int64) (BaseMessage, error) {
	if apiKey == "" {
		apiKey = os.Getenv("BLOCKNATIVE_DAPP_ID")
	}
	net, err := NetName(netID)
	if err != nil {
		return BaseMessage{}, err
	}
	return BaseMessage{
		Timestamp: time.Now(),
		DappID:    apiKey,
		Blockchain: Blockchain{
			System:  "ethereum",
			Network: net,
		},
	}, nil
}
