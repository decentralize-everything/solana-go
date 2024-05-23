package rpc

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func (cl *Client) SendBundle(
	ctx context.Context,
	transaction *solana.Transaction,
) (signature solana.Signature, err error) {
	txData, err := transaction.MarshalBinary()
	if err != nil {
		return solana.Signature{}, fmt.Errorf("encode transaction failed: %v", err)
	}
	signature = transaction.Signatures[0]

	var sig string
	err = cl.rpcClient.CallForInto(ctx, &sig, "sendBundle", []interface{}{
		[]string{base58.Encode(txData)},
	})
	return
}
