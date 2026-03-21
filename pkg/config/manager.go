package config

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/lattesec/log"
)

type ConfigManager struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	cancel context.CancelFunc

	cfg        atomic.Pointer[Config]
	customPath atomic.Value
}

func NewConfigManager(customPath string) (*ConfigManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	cm := &ConfigManager{
		ctx:    ctx,
		wg:     &wg,
		cancel: cancel,
	}
	cm.customPath.Store(customPath)

	if err := cm.Load(); err != nil {
		return nil, err
	}

	wg.Go(cm.watchConfig)
	return cm, nil
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.cfg.Load()
}

func (cm *ConfigManager) SetCustomPath(path *string) {
	cm.customPath.Store(path)
}

func (cm *ConfigManager) Close() {
	cm.cancel()
	cm.wg.Wait()
}

func (cm *ConfigManager) Load() error {
	cfg, _, err := LoadConfig(cm.customPath.Load().(string))
	if err != nil {
		return err
	}
	cm.cfg.Store(cfg)
	return nil
}

// Watches SIGHUP
func (cm *ConfigManager) watchConfig() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP)

	for {
		select {
		case <-cm.ctx.Done():
			log.Info().WithMeta("scope", "cfg").Msg("shutting down config watcher").Send()
			return
		case <-sigCh:
			if err := cm.Load(); err != nil {
				log.Error().WithMeta("scope", "cfg").Msgf("failed to reload config: %v", err).Send()
			}
		}
	}
}
