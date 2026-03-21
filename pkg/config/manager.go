package config

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/lattesec/log"
)

type ConfigManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg        atomic.Pointer[Config]
	customPath atomic.Value
}

func NewConfigManager(customPath string) *ConfigManager {
	ctx, cancel := context.WithCancel(context.Background())

	cm := &ConfigManager{
		ctx:    ctx,
		cancel: cancel,
	}
	cm.customPath.Store(customPath)

	go cm.watchConfig()
	return cm
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.cfg.Load()
}

func (cm *ConfigManager) SetCustomPath(path *string) {
	cm.customPath.Store(path)
}

func (cm *ConfigManager) Close() {
	cm.cancel()
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
			return
		case <-sigCh:
			if err := cm.Load(); err != nil {
				log.Error().WithMeta("scope", "cfg").Msgf("failed to reload config: %v", err).Send()
			}
		}
	}
}
