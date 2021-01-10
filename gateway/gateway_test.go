package gateway

import (
	"chain_proto/testutils/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestGateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	bc := mocks.NewMockBlockchain(ctrl)
	g := New(bc)

	if err := g.Start(); err != nil {
		t.Error("failed to start the gateway", err)
	}

	time.Sleep(time.Hour)
}
