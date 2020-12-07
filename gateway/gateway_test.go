package gateway

import "testing"

func TestHTTPServer(t *testing.T) {
	g := New()
	if err := g.Start(); err != nil {
		t.Error(err)
	}
}
