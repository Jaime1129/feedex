package jobs

import (
	"context"
	"log"
	"time"

	"github.com/jaime1129/fedex/internal/components"
)

type DataTracker interface {
}

type dataTracker struct {
	ethScanCli components.EthScanCli
}

func (t *dataTracker) TrackLiveData(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	select {
	case <-ticker.C:
		log.Println("refresh live data...")
		t.trackLiveData(ctx)
	case <-ctx.Done():
		log.Println("live data tracker stopped")
	}
}

func (t *dataTracker) trackLiveData(ctx context.Context) {
	t.ethScanCli.QueryBlock()
}

func (t *dataTracker) TrackHistoricalData(ctx context.Context) {

}
