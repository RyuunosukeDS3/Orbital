package test

import (
	"sync/atomic"
	"testing"
	"time"

	"orbital/internal/argocd"
	"orbital/internal/monitoring"
	"orbital/internal/valve"

	"github.com/stretchr/testify/assert"
)

func TestPlayersOfflineScalesDownAfter5EmptyQueries(t *testing.T) {
	originalQuery := valve.QueryVrisingPlayersNum
	originalSetReplicaCount := argocd.SetReplicaCount
	originalNewTicker := monitoring.NewTicker

	defer func() {
		valve.QueryVrisingPlayersNum = originalQuery
		argocd.SetReplicaCount = originalSetReplicaCount
		monitoring.NewTicker = originalNewTicker
	}()

	monitoring.NewTicker = func(d time.Duration) *time.Ticker {
		return time.NewTicker(10 * time.Millisecond)
	}

	var callCounter int32
	valve.QueryVrisingPlayersNum = func(ip, port string) (uint8, error) {
		atomic.AddInt32(&callCounter, 1)
		return 0, nil
	}	

	var called int32
	argocd.SetReplicaCount = func(appName, count string) error {
		atomic.StoreInt32(&called, 1)
		assert.Equal(t, "test-app", appName)
		assert.Equal(t, "0", count)
		return nil
	}

	go monitoring.PlayersOffline("127.0.0.1", "9876", "test-app")

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(1), atomic.LoadInt32(&called), "expected SetReplicaCount to be called after 5 empty ticks")
}
