package rpc_test

import (
	"context"
	"testing"

	"github.com/decentralize-everything/solana-go/rpc"
	"github.com/gagliardetto/solana-go"
)

func TestGetProgramAccountsV2(t *testing.T) {
	client := rpc.New("https://mainnet.helius-rpc.com/?api-key=8ce36b60-ec4f-4e8e-afe7-9692a6a469d3")

	result, err := client.GetProgramAccountsV2(
		context.TODO(),
		solana.MustPublicKeyFromBase58("AddressLookupTab1e1111111111111111111111111"),
		&rpc.GetProgramAccountsV2Opts{
			Limit: 10,
			Filters: []rpc.RPCFilter{
				{DataSize: uint64(88 + 32)},
			},
		},
	)
	if err != nil {
		t.Fatalf("GetProgramAccountsV2 failed: %v", err)
	}

	t.Logf("GetProgramAccountsV2 returned %d accounts", len(result.Accounts))
}
