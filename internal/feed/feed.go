// Package feed handles the storage and intricacies of
// having a minimal loss camera server.
package feed

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/caffeine-addictt/camserver/pkg/config"
	"github.com/lattesec/log"
)

type FeedManager struct {
	wg      *sync.WaitGroup
	ctx     *context.Context
	ctxDone context.CancelFunc

	cameras    []*config.CameraCfg
	archiveDir atomic.Value
}

// NewFeedManager follows the RAII model
func NewFeedManager(ctx context.Context, rootWg *sync.WaitGroup, cameras ...*config.CameraCfg) *FeedManager {
	fmCtx, fmCtxDone := context.WithCancel(context.Background())

	fm := &FeedManager{
		wg:      &sync.WaitGroup{},
		ctx:     &fmCtx,
		ctxDone: fmCtxDone,
		cameras: cameras,
	}

	rootWg.Go(func() {
		<-ctx.Done()
		fm.Stop()
		log.Info().Msg("shutting down feed manager").Send()
		fmt.Println("shutting down feed manager")
	})

	return fm
}

func (fm *FeedManager) Stop() {
	fm.ctxDone()
	fm.wg.Wait()
}
