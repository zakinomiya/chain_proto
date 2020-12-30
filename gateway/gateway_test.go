package gateway

import (
	"chain_proto/config"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	g := New(&config.Config.Network)
	if err := g.Start(); err != nil {
		t.Error(err)
	}
}
