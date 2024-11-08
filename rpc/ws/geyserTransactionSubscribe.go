package ws

import (
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"
)

type TransactionResult struct {
	Transaction struct {
		Transaction []string `json:"transaction"`
		Meta        struct {
			Err               interface{}            `json:"err"`
			Status            interface{}            `json:"status"`
			Fee               uint64                 `json:"fee"`
			PreBalances       []uint64               `json:"preBalances"`
			PostBalances      []uint64               `json:"postBalances"`
			PreTokenBalances  []rpc.TokenBalance     `json:"preTokenBalances"`
			PostTokenBalances []rpc.TokenBalance     `json:"postTokenBalances"`
			LogMessages       []string               `json:"logMessages"`
			LoadedAddresses   rpc.LoadedAddresses    `json:"loadedAddresses"`
			InnerInstructions []rpc.InnerInstruction `json:"innerInstructions"`
		} `json:"meta"`
	} `json:"transaction"`
	Signature solana.Signature `json:"signature"`
}

func (cl *Client) SupportGeyser() bool {
	return strings.Contains(cl.rpcURL, "helius")
}

func (cl *Client) GeyserTransactionSubscribe(addresses []string) (*LogSubscription, error) {
	genSub, err := cl.subscribe(
		[]interface{}{
			map[string]interface{}{
				"vote":           false,
				"failed":         false,
				"accountInclude": addresses,
			},
		},
		map[string]interface{}{
			"commitment":                     "processed",
			"encoding":                       "base64",
			"transactionDetails":             "full",
			"showRewards":                    true,
			"maxSupportedTransactionVersion": 0,
		},
		"transactionSubscribe",
		"transactionUnsubscribe",
		func(msg []byte) (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					zlog.Error("decoderFunc", zap.Any("recover", r))
				}
			}()

			var res TransactionResult
			err := decodeResponseFromMessage(msg, &res)
			if err != nil {
				return nil, err
			}

			var logResult LogResult
			logResult.Value.Err = res.Transaction.Meta.Err
			logResult.Value.Logs = res.Transaction.Meta.LogMessages
			logResult.Value.Signature = res.Signature
			logResult.Value.TransactionResult = &res
			return &logResult, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &LogSubscription{
		sub: genSub,
	}, nil
}
