package rpc

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
)

func TestSendBundle(t *testing.T) {
	client := New("https://buckeroos-megalopses-tkoouqmwbw-dedicated.helius-rpc.com?api-key=001eb0e9-7c61-4567-87f8-bcdaf22c9589")
	signature, err := client.SendBundle(
		context.TODO(),
		&solana.Transaction{Signatures: []solana.Signature{{}}},
	)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(signature)
}
