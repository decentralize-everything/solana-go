package rpc

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

func (cl *Client) GetProgramAccountsV2(
	ctx context.Context,
	publicKey solana.PublicKey,
	opts *GetProgramAccountsV2Opts,
) (out GetProgramAccountsV2Result, err error) {
	obj := M{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = string(opts.Commitment)
		}
		if len(opts.Filters) != 0 {
			obj["filters"] = opts.Filters
		}
		if opts.Encoding != "" {
			obj["encoding"] = opts.Encoding
		}
		if opts.DataSlice != nil {
			obj["dataSlice"] = M{
				"offset": opts.DataSlice.Offset,
				"length": opts.DataSlice.Length,
			}
		}
		if opts.Limit != 0 {
			obj["limit"] = opts.Limit
		}
		if opts.PaginationKey != "" {
			obj["paginationKey"] = opts.PaginationKey
		}
	}

	params := []interface{}{publicKey, obj}

	err = cl.rpcClient.CallForInto(ctx, &out, "getProgramAccountsV2", params)
	return
}
