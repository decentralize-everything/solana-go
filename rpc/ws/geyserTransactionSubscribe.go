package ws

import (
	"strings"

	"github.com/gagliardetto/solana-go"
	"go.uber.org/zap"
)

type TransactionResult struct {
	Transaction struct {
		Transaction []string `json:"transaction"`
		Meta        struct {
			Log []string `json:"logMessages"`
		} `json:"meta"`
	} `json:"transaction"`
	Signature solana.Signature `json:"signature"`
}

func (cl *Client) SupportGeyser() bool {
	return strings.Contains(cl.rpcURL, "helius")
}

func (cl *Client) GeyserTransactionSubscribe(address string) (*LogSubscription, error) {
	genSub, err := cl.subscribe(
		[]interface{}{
			map[string]interface{}{
				"vote":   false,
				"failed": false,
				"accountInclude": []string{
					address,
				},
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
			logResult.Value.Logs = res.Transaction.Meta.Log
			logResult.Value.Signature = res.Signature
			logResult.Value.TransactionBase64 = res.Transaction.Transaction[0]
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
