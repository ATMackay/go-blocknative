package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ATMackay/go-blocknative/client"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

const (
	netName      = "main"
	contractAddr = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	methodName   = "transfer"
	UsdtABI      = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_upgradedAddress","type":"address"}],"name":"deprecate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"deprecated","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_evilUser","type":"address"}],"name":"addBlackList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"upgradedAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balances","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"maximumFee","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"_totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"unpause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_maker","type":"address"}],"name":"getBlackListStatus","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"allowed","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"paused","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"who","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"pause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getOwner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newBasisPoints","type":"uint256"},{"name":"newMaxFee","type":"uint256"}],"name":"setParams","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"}],"name":"issue","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"}],"name":"redeem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"basisPointsRate","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"isBlackListed","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_clearedUser","type":"address"}],"name":"removeBlackList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"MAX_UINT","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_blackListedUser","type":"address"}],"name":"destroyBlackFunds","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"_initialSupply","type":"uint256"},{"name":"_name","type":"string"},{"name":"_symbol","type":"string"},{"name":"_decimals","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Issue","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Redeem","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAddress","type":"address"}],"name":"Deprecate","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"feeBasisPoints","type":"uint256"},{"indexed":false,"name":"maxFee","type":"uint256"}],"name":"Params","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_blackListedUser","type":"address"},{"indexed":false,"name":"_balance","type":"uint256"}],"name":"DestroyedBlackFunds","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_user","type":"address"}],"name":"AddedBlackList","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_user","type":"address"}],"name":"RemovedBlackList","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[],"name":"Pause","type":"event"},{"anonymous":false,"inputs":[],"name":"Unpause","type":"event"}]`
)

func main() {

	mempMon, err := client.New(context.Background(), client.Opts{
		Scheme: "wss",
		Host:   "api.blocknative.com",
		Path:   "/v0",
		APIKey: os.Getenv("BLOCKNATIVE_DAPP_ID"),
	})

	exitOnErr(err, "create blocknative client")

	baseMsg := client.BaseMessage{
		Timestamp: time.Now(),
		DappID:    os.Getenv("BLOCKNATIVE_DAPP_ID"),
		Blockchain: client.Blockchain{
			System:  "ethereum",
			Network: netName,
		},
	}

	exitOnErr(mempMon.Initialize(baseMsg), "initialize subs")

	var abi interface{}
	exitOnErr(json.Unmarshal([]byte(UsdtABI), &abi), "marshal abi")

	cfgMsg := client.NewConfig(
		contractAddr,
		true,
		abi,
	)
	cfgMsg.Filters = []map[string]string{
		{
			"contractCall.methodName": methodName,
			"_propertySearch":         "true",
		},
	}

	cfgMsgWithBase := client.NewConfiguration(baseMsg, cfgMsg)

	msg, err := json.MarshalIndent(cfgMsgWithBase, "", "  ")
	exitOnErr(err, "config message marshal")
	log.Println("cfgMsgWithBase", string(msg))

	exitOnErr(mempMon.NewEventSubscription(cfgMsgWithBase), "config subs")
	log.Print("subscription created   ", "network:", netName, "   contract:", contractAddr, "    method:", methodName)

	evSub := mempMon.SubscriptionRegistry()[cfgMsg.Scope]
	eventChan := evSub.Events()

	type txIndex struct {
		blockNumber      int
		transactionIndex int
		call             string
	}

	indexer := make(chan txIndex)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var counter int
		for {
			time.Sleep(2 * time.Second)
			if counter > 5 {
				return
			}
			counter++
			if len(eventChan) < 1 {
				log.Printf("channel empty\n")
				continue
			}
			e, ok := <-eventChan
			if !ok {
				break
			}
			ev := e.(client.EthTxPayload)
			if len(ev.Event.Transaction.Input) < 10 {
				continue
			}
			jev, _ := json.MarshalIndent(ev, "", "  ")
			log.Printf("msg: %+v \n", string(jev))
			s, ok := parseInput(ev.Event.Transaction.Input).(string)
			if !ok {
				s = "could not parse data"
			}
			indexer <- txIndex{blockNumber: ev.Event.Transaction.BlockNumber, transactionIndex: ev.Event.Transaction.TransactionIndex, call: s}
		}
		return
	}()

	wg.Wait()
	evSub.Unsubscribe()
	for len(indexer) != 0 {
		tx, ok := <-indexer
		if !ok {
			break
		}
		log.Printf("transaction: %v\n", tx)
	}
	time.Sleep(2 * time.Second)
	err = <-evSub.Err()
	log.Printf("subscription close message: %v\n", err)
}

func exitOnErr(err error, msg string) {
	logger := log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
	if err != nil {
		logger.Output(2, fmt.Sprintf("root execution error:%+v msg:%+v", err, msg))
		os.Exit(1)
	}
}

func parseInput(input string) interface{} {
	abiT, err := abi.JSON(strings.NewReader(UsdtABI))
	exitOnErr(err, "loading the abi")

	inputData, err := hex.DecodeString(input[10:])
	exitOnErr(err, "input decode")

	method, exist := abiT.Methods[methodName]
	if !exist {
		exitOnErr(errors.New("method doesn't exists in the abi"), "")
	}

	output, err := method.Inputs.Unpack(inputData)
	exitOnErr(err, "args unpack")

	return output
}
