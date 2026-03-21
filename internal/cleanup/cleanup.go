// Package cleanup
package cleanup

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lattesec/log"
)

var (
	cleanupStack      []func() error
	cleanupRwLock     sync.RWMutex
	errorCleanupStack []func() error
	errorRwLock       sync.RWMutex
)

// Schedule adds a function to the cleanup stack
//
// Cleanups are ran no matter what, before error cleanup
func Schedule(fn func() error) {
	cleanupRwLock.Lock()
	defer cleanupRwLock.Unlock()
	cleanupStack = append(cleanupStack, fn)
}

// ScheduleError adds a function to the error cleanup stack
//
// Error cleanups are ran ONLY when an error occurs, and
// are ran after all regular cleanups
func ScheduleError(fn func() error) {
	errorRwLock.Lock()
	defer errorRwLock.Unlock()
	errorCleanupStack = append(errorCleanupStack, fn)
}

// Cleanup handles cleanup functions
func Cleanup() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("recovered from panic: %v\n", r).Send()
		}
	}()

	cleanupRwLock.Lock()
	defer cleanupRwLock.Unlock()

	log.Debug().Msgf("cleaning up %d items...\n", len(cleanupStack)).Send()
	for i := len(cleanupStack) - 1; i >= 0; i-- {
		if err := cleanupStack[i](); err != nil {
			log.Error().Msgf("error while cleaning up: %v\n", err).Send()
		}
	}

	cleanupStack = []func() error{}
}

// CleanupError handles all error cleanup functions
func CleanupError() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("recovered from panic: %v\n", r).Send()
		}
	}()

	errorRwLock.Lock()
	defer errorRwLock.Unlock()

	log.Debug().Msgf("cleaning up %d items due to error...\n", len(errorCleanupStack)).Send()
	for i := len(errorCleanupStack) - 1; i >= 0; i-- {
		if err := errorCleanupStack[i](); err != nil {
			log.Error().Msgf("error while cleaning up: %v\n", err).Send()
		}
	}

	errorCleanupStack = []func() error{}
}

func Watch(rootCtx context.Context, rootDone context.CancelFunc) *sync.WaitGroup {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	wg := WatchContext(rootCtx, rootDone, ctx)
	return wg
}

func WatchContext(rootCtx context.Context, rootDone context.CancelFunc, sigCtx context.Context) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Go(func() {
		select {
		case <-rootCtx.Done():
		case <-sigCtx.Done():
			log.Debug().Msg("cleaning up...").Send()

			rootDone()
			Cleanup()
			CleanupError()
			log.Sync()
		}
	})

	return &wg
}
