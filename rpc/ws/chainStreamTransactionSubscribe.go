package ws

import (
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type ChainStreamTransactionResult struct {
	Context struct {
		SlotStatus string           `json:"slotStatus"`
		NodeTime   time.Time        `json:"nodeTime"`
		IsVote     bool             `json:"isVote"`
		Signature  solana.Signature `json:"signature"`
		Index      int              `json:"index"`
	} `json:"context"`
	Value struct {
		BlockTime time.Time `json:"blockTime"`
		Meta      struct {
			Err               interface{}            `json:"err"`
			Fee               uint64                 `json:"fee"`
			InnerInstructions []rpc.InnerInstruction `json:"innerInstructions"`
			LoadedAddresses   rpc.LoadedAddresses    `json:"loadedAddresses"`
			LogMessages       []string               `json:"logMessages"`
			PostBalances      []uint64               `json:"postBalances"`
			PostTokenBalances []rpc.TokenBalance     `json:"postTokenBalances"`
			PreBalances       []uint64               `json:"preBalances"`
			PreTokenBalances  []rpc.TokenBalance     `json:"preTokenBalances"`
		} `json:"meta"`
		Slot        uint64             `json:"slot"`
		Transaction solana.Transaction `json:"transaction"`
	} `json:"value"`
}

func (cl *Client) SupportChainStream() bool {
	return strings.Contains(cl.rpcURL, "syndica")
}

func (cl *Client) ChainStreamTransactionSubscribe(addresses []string) (*LogSubscription, error) {
	genSub, err := cl.subscribe(
		nil,
		map[string]interface{}{
			"network":  "solana-mainnet",
			"verified": false,
			"filter": map[string]interface{}{
				"excludeVotes": true,
				"commitment":   "processed",
				"accountKeys": map[string]interface{}{
					"oneOf": addresses,
				},
			},
		},
		"transactionsSubscribe",
		"transactionsUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res ChainStreamTransactionResult
			err := decodeResponseFromMessage(msg, &res)
			if err != nil {
				return nil, err
			}

			if res.Value.Meta.Err != nil {
				return nil, nil
			}
			transaction := res.Value.Transaction.MustToBase64()

			var logResult LogResult
			logResult.Value.Err = res.Value.Meta.Err
			logResult.Value.Logs = res.Value.Meta.LogMessages
			logResult.Value.Signature = res.Context.Signature
			logResult.Value.TransactionResult = &TransactionResult{
				Transaction: struct {
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
				}{
					Transaction: []string{
						transaction,
					},
					Meta: struct {
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
					}{
						Err:               nil,
						Status:            nil,
						Fee:               res.Value.Meta.Fee,
						PreBalances:       res.Value.Meta.PreBalances,
						PostBalances:      res.Value.Meta.PostBalances,
						PreTokenBalances:  res.Value.Meta.PreTokenBalances,
						PostTokenBalances: res.Value.Meta.PostTokenBalances,
						LogMessages:       res.Value.Meta.LogMessages,
						LoadedAddresses:   res.Value.Meta.LoadedAddresses,
						InnerInstructions: res.Value.Meta.InnerInstructions,
					},
				},
				Signature: res.Context.Signature,
			}

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
