package ws

import (
	"context"
	"testing"
)

func BenchmarkChainStreamTransactionSubscribe(t *testing.B) {
	cl, err := Connect(context.Background(), "wss://chainstream.api.syndica.io/api-token/27e1wWCAbWAy1EJwEb3U2dvTeee6TRgk76eEjysHdgcNt8VWUfHB6kJtHm1Yhou5MxiiT67JHqYhHTQFhn4djymENFm6By7cXQwNmaS5c9YyHzuB7FZk6qmV61VBX43HY8JTF1XpdgFLSANPmm2QY9eoqaj25fP2z6uz39CPCTeoH5bV2gim5n8YEUCKpweFLayGPyEjCyUTiM67jr29Mo935DBKqWfB9UFzv2WaKGQoXGod6gPzwHFtWDhsK6tgYTxD2nt3B3eaiNUHcU5811xzr3AL6UWCVS8LUnQKf6JTA7fCeed1ZZgVjcbuEEZh6etk5UjKmUUSiF2kVPDJi5kfpuhFF12j5BmLtmg3ixwvhLF6ydGg7MqooTbTgpmvL9f4W6F4yF5Sj3oFep7hrtq7aBd1dg1FG48HgHKuJtoNdUcKJ3RmHQidwjF4e38ygaFtgbYJj6Rc6U3T3xj3DTQtGyeY5G4UBda5b6EKD4brzPkDWMQ1D2xBtDVGi")
	if err != nil {
		t.Log(err)
		return
	}

	sub, err := cl.ChainStreamTransactionSubscribe([]string{"675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"})
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