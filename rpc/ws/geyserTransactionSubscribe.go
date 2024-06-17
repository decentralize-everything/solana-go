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
			Fee               uint64             `json:"fee"`
			PreBalances       []uint64           `json:"preBalances"`
			PostBalances      []uint64           `json:"postBalances"`
			PreTokenBalances  []rpc.TokenBalance `json:"preTokenBalances"`
			PostTokenBalances []rpc.TokenBalance `json:"postTokenBalances"`
			LogMessages       []string           `json:"logMessages"`
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

			tx, err := solana.TransactionFromBase64(res.Transaction.Transaction[0])
			if err != nil {
				return nil, err
			}

			var logResult LogResult
			logResult.Value.Logs = res.Transaction.Meta.LogMessages
			logResult.Value.Signature = res.Signature
			logResult.Value.TransactionMeta = &rpc.TransactionMeta{
				Fee:               res.Transaction.Meta.Fee,
				PreBalances:       res.Transaction.Meta.PreBalances,
				PostBalances:      res.Transaction.Meta.PostBalances,
				PreTokenBalances:  res.Transaction.Meta.PreTokenBalances,
				PostTokenBalances: res.Transaction.Meta.PostTokenBalances,
			}
			logResult.Value.Message.AccountKeys = tx.Message.AccountKeys
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
