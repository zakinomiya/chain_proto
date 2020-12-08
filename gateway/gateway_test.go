package gateway

import (
	"go_chain/config"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	g := New(&config.Config.Network)
	if err := g.Start(); err != nil {
		t.Error(err)
	}
}
