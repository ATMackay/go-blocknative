package client

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/params"
)

func ParseGas(msg *EthTxPayload) (gasBaseFeeGwei, gasTipGwei float64, err error) {
	gasBaseFee, err := strconv.ParseFloat(msg.Event.Transaction.MaxFeePerGas, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing gas base fee: %v, err=%v", msg.Event.Transaction.MaxFeePerGas, err)
	}

	gasBaseFeeGwei = gasBaseFee / params.GWei

	gasTip, err := strconv.ParseFloat(msg.Event.Transaction.MaxPriorityFeePerGas, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing max priority fee: %v, err=%v", msg.Event.Transaction.MaxPriorityFeePerGas, err)
	}
	gasTipGwei = gasTip / params.GWei

	return gasBaseFeeGwei, gasTipGwei, nil
}
