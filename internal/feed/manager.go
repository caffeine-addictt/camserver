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

	watchers   []*Watcher
	archiveDir atomic.Value
}

// NewFeedManager follows the RAII model
func NewFeedManager(ctx context.Context, rootWg *sync.WaitGroup, dir string, cameras ...config.CameraCfg) *FeedManager {
	fmCtx, fmCtxDone := context.WithCancel(context.Background())

	fm := &FeedManager{
		wg:         &sync.WaitGroup{},
		ctx:        &fmCtx,
		ctxDone:    fmCtxDone,
		watchers:   []*Watcher{},
		archiveDir: atomic.Value{},
	}

	rootWg.Go(func() {
		<-ctx.Done()
		fm.Stop()
		log.Info().Msg("shutting down feed manager").Send()
		fmt.Println("shutting down feed manager")
	})

	fm.archiveDir.Store(dir)
	fm.UpdateCameras(cameras...)
	return fm
}

func (fm *FeedManager) Stop() {
	fm.ctxDone()
	fm.wg.Wait()
}

func (fm *FeedManager) UpdateCameras(cameras ...config.CameraCfg) {
	newCameras := make(map[string]*config.CameraCfg, len(cameras))
	for _, w := range cameras {
		newCameras[fmt.Sprintf("%s-%s", w.Name, w.Rtsp.String())] = &w
	}

	dir := fm.archiveDir.Load().(string)

	for _, w := range fm.watchers {
		key := fmt.Sprintf("%s-%s", w.Camera.Name, w.Camera.Rtsp.String())

		if dir != w.ArchiveDir {
			w.Stop()
			continue
		}

		if _, ok := newCameras[key]; !ok {
			w.Stop()
			continue
		}

		delete(newCameras, key)
	}

	for _, w := range newCameras {
		newW := NewWatcher(*fm.ctx, fm.wg, w, dir)
		newW.Start()
		fm.watchers = append(fm.watchers, newW)
	}
}

func (fm *FeedManager) UpdateArchiveDir(archiveDir string) {
	fm.archiveDir.Store(archiveDir)
}

func (fm *FeedManager) ArchiveDir() string {
	return fm.archiveDir.Load().(string)
}
