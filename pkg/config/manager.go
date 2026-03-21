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

type (
	ConfigManagerCallback func(newCfg, oldCfg *Config)

	ConfigManager struct {
		ctx    context.Context
		wg     *sync.WaitGroup
		cancel context.CancelFunc

		cfg        atomic.Pointer[Config]
		customPath atomic.Value

		callbacks []ConfigManagerCallback
	}
)

func NewConfigManager(customPath string) (*ConfigManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	cm := &ConfigManager{
		ctx:       ctx,
		wg:        &wg,
		cancel:    cancel,
		callbacks: []ConfigManagerCallback{},
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

// RegisterCallback adds a new func that will be called when config changes
func (cm *ConfigManager) RegisterCallback(cb ConfigManagerCallback) {
	cm.callbacks = append(cm.callbacks, cb)
}

func (cm *ConfigManager) load() (*Config, error) {
	cfg, _, err := LoadConfig(cm.customPath.Load().(string))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cm *ConfigManager) Load() error {
	cfg, err := cm.load()
	if err != nil {
		return err
	}
	cm.cfg.Store(cfg)
	return nil
}

// Watches SIGUSR1
func (cm *ConfigManager) watchConfig() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGUSR1)

	for {
		select {
		case <-cm.ctx.Done():
			log.Info().WithMeta("scope", "cfg").Msg("shutting down config watcher").Send()
			return
		case <-sigCh:
			cfg, err := cm.load()
			if err != nil {
				log.Error().WithMeta("scope", "cfg").Msgf("failed to reload config: %v", err).Send()
				continue
			}

			log.Info().WithMeta("scope", "cfg").Msg("SIGUSR1, reloded config").Send()
			old := cm.cfg.Swap(cfg)
			for _, cb := range cm.callbacks {
				cb(old, cfg)
			}
		}
	}
}
