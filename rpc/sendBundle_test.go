package rpc

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
)

func TestSendBundle(t *testing.T) {
	client := New("https://mainnet.block-engine.jito.wtf/api/v1/bundles")
	signature, err := client.SendBundle(
		context.TODO(),
		&solana.Transaction{},
	)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(signature)
}
