package ws

import (
	"context"
	"testing"
)

func BenchmarkGeyserTransactionSubscribe(t *testing.B) {
	cl, err := Connect(context.Background(), "wss://atlas-mainnet.helius-rpc.com/?api-key=f3c02a35-f460-48c5-95ac-392942fc621d")
	if err != nil {
		t.Log(err)
		return
	}

	sub, err := cl.GeyserTransactionSubscribe("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	if err != nil {
		t.Log(err)
		return
	}

	count := 0
	for {
		event := <-sub.Response()
		count += 1
		if count%100 == 0 {
			t.Log(count)
		}
		if event.Value.Err == nil {
			t.Log(event.Value.Signature)
		}
	}
}
